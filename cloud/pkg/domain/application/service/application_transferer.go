// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	applicationentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"time"
)

type statusHandler func(appDeployInfo *applicationentity.ApplicationDeployInfo)

// status transfer
var statusHandlerMap = map[string]statusHandler{
	constants.StatusInit:      handleStatusInit,
	constants.StatusDeploying: handleStatusDeploying,
	constants.StatusRunning:   handleStatusRunning,
	constants.StatusDeleting:  handleStatusDeleting,
}

// CreateNodeCheckLoop transferer of edge application status
func CreateApplicationTickerLoop(ctx context.Context, config *configs.EdgeCloudConfig) {
	duration := time.Duration(config.Portal.ApplicationCheckerInterval) * time.Second
	checkTicker := time.NewTimer(duration)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("CreateApplicationTickerLoop was terminated")
			return
		case <-checkTicker.C:
			ApplicationTicker(config)
			checkTicker.Reset(duration)
		}
	}
}

func ApplicationTicker(config *configs.EdgeCloudConfig) {
	var appDeployInfos []applicationentity.ApplicationDeployInfo
	var totalTmp int64
	var total int
	offset := 0

	apulisdb.Db.Model(&applicationentity.ApplicationDeployInfo{}).Count(&totalTmp)
	total = int(totalTmp)

	logger.Debugf("ApplicationTicker total application count = %d", total)
	if total == 0 {
		return
	} else if total < constants.TransferCountEach {
		total = constants.TransferCountEach
	}

	for total >= constants.TransferCountEach {
		res := apulisdb.Db.Offset(offset).Limit(constants.TransferCountEach).Find(&appDeployInfos)
		if res.Error != nil {
			logger.Errorf("query application deploy failed. err = %v", res.Error)
		} else {
			for i := 0; i < int(res.RowsAffected); i++ {
				logger.Debugf("ApplicationTicker handle application = %v", appDeployInfos[i])
				statusHandlerMap[appDeployInfos[i].Status](&appDeployInfos[i])
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

func handleStatusInit(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	// first: update status
	appDeployInfo.Status = constants.StatusDeploying
	appDeployInfo.UpdateAt = time.Now()
	err := applicationentity.UpdateAppDeploy(appDeployInfo)
	if err != nil {
		logger.Infof("update deployment failed when deploying!")
		return
	}

	// second: deploy to k8s
	err = CreateK8sPod(appDeployInfo)
	if err != nil {
		logger.Infof("create pod failed! err = %v", err)
		return
	}

	logger.Infof("create pod succ! status to %s", constants.StatusDeploying)
}

func handleStatusDeploying(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	// check pod status, if ok, transfer to StatusRunning; if failed, retry
	pod, err := GetK8sPod(appDeployInfo)
	if err != nil { // if failed, try next time
		return
	}

	if pod.Status.Phase == corev1.PodRunning {
		appDeployInfo.Status = constants.StatusRunning
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
		}
	}

	// TODO other status handle
}

func handleStatusRunning(appDeployInfo *applicationentity.ApplicationDeployInfo) {

}

func handleStatusDeleting(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	// first: undeploy to k8s
	_, err := GetK8sPod(appDeployInfo)
	if err != nil { // if failed, try next time
		// metav1.StatusReasonNotFound
		logger.Infof("handleStatusDeleting, GetK8sPod err, err = %v", err)
		return
	}

	err = DeleteK8sPod(appDeployInfo)
	if err != nil { // if failed, try next time
		return
	}

	// second: db record delete
	err = applicationentity.DeleteAppDeploy(appDeployInfo)
	if err == nil {
		logger.Infof("handleStatusDeleting, DeleteAppDeploy err, err = %v", err)
	}
}

func CreateK8sPod(dbInfo *applicationentity.ApplicationDeployInfo) error {
	podClient, err := utils.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return err
	}

	verInfo, err := applicationentity.GetApplicationVersion(dbInfo.ClusterId, dbInfo.GroupId, dbInfo.UserId, dbInfo.AppName, dbInfo.Version)
	if err != nil {
		return err
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName(dbInfo.ClusterId, dbInfo.GroupId, dbInfo.UserId, dbInfo.AppName, dbInfo.Version, dbInfo.NodeName),
		},
		Spec: corev1.PodSpec{
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": dbInfo.NodeName,
			},
			Containers: []corev1.Container{
				{
					Name:  verInfo.ContainerImage,
					Image: verInfo.ContainerImagePath,
				},
			},
		},
	}

	_, err = podClient.Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetK8sPod(dbInfo *applicationentity.ApplicationDeployInfo) (*corev1.Pod, error) {
	podClient, err := utils.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return nil, err
	}

	pod, err := podClient.Get(context.Background(), podName(dbInfo.ClusterId, dbInfo.GroupId, dbInfo.UserId, dbInfo.AppName, dbInfo.Version, dbInfo.NodeName), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func DeleteK8sPod(dbInfo *applicationentity.ApplicationDeployInfo) error {
	podClient, err := utils.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return err
	}

	return podClient.Delete(context.Background(), podName(dbInfo.ClusterId, dbInfo.GroupId, dbInfo.UserId, dbInfo.AppName, dbInfo.Version, dbInfo.NodeName), metav1.DeleteOptions{})
}

func podName(clusterId int64, groupId int64, userId int64, appName string, version string, nodeName string) string {
	return appName + "-" + strconv.FormatInt(userId, 10) + "-" + "deployment"
}
