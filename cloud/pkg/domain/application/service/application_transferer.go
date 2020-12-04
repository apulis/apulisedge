// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	applicationentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/utils"
	v1 "k8s.io/api/apps/v1"
	k8scorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
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
	err = CreateK8sDeployment(appDeployInfo)
	if err != nil {
		logger.Infof("create deployment failed! err = %v", err)
		return
	}

	logger.Infof("create deployment succ! status to %s", constants.StatusDeploying)
}

func handleStatusDeploying(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	// check deployment status, if ok, transfer to StatusRunning; if failed, retry
	deploy, err := GetK8sDeployment(appDeployInfo)
	if err != nil { // if failed, try next time
		return
	}

	// TODO add replicas to db info
	if deploy.Status.ReadyReplicas == 1 {
		appDeployInfo.Status = constants.StatusRunning
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
		}
	}
}

func handleStatusRunning(appDeployInfo *applicationentity.ApplicationDeployInfo) {

}

func handleStatusDeleting(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	// first: undeploy to k8s
	_, err := GetK8sDeployment(appDeployInfo)
	if err != nil { // if failed, try next time
		// metav1.StatusReasonNotFound
		logger.Infof("handleStatusDeleting, GetK8sDeployment err, err = %v", err)
		return
	}

	err = DeleteK8sDeployment(appDeployInfo)
	if err != nil { // if failed, try next time
		return
	}

	// second: db record delete
	err = applicationentity.DeleteAppDeploy(appDeployInfo)
	if err == nil {
		logger.Infof("handleStatusDeleting, DeleteAppDeploy err, err = %v", err)
	}
}

func deploymentName(appName string, userId int64) string {
	return appName + "-" + strconv.FormatInt(userId, 10) + "-" + "deployment"
}

func CreateK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) error {
	deployClient, err := utils.GetDeploymentClient(dbInfo.Namespace)
	if err != nil {
		return err
	}

	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName(dbInfo.AppName, dbInfo.UserId), // deployment名字
		},
		Spec: v1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1), // 指定副本数
			Selector: &metav1.LabelSelector{ // 指定标签
				MatchLabels: map[string]string{
					"app": dbInfo.AppName,
				},
			},
			Template: k8scorev1.PodTemplateSpec{ // 容器模板
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": dbInfo.AppName,
					},
				},
				Spec: k8scorev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": dbInfo.NodeName,
					},
					Containers: []k8scorev1.Container{
						{
							Name:  dbInfo.ContainerImage,
							Image: dbInfo.ContainerImagePath,
							Ports: []k8scorev1.ContainerPort{
								{
									ContainerPort: int32(dbInfo.ContainerPort),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = deployClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) (*v1.Deployment, error) {
	deployClient, err := utils.GetDeploymentClient(dbInfo.Namespace)
	if err != nil {
		return nil, err
	}

	deploy, err := deployClient.Get(context.Background(), deploymentName(dbInfo.AppName, dbInfo.UserId), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

func DeleteK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) error {
	deployClient, err := utils.GetDeploymentClient(dbInfo.Namespace)
	if err != nil {
		return err
	}

	return deployClient.Delete(context.Background(), deploymentName(dbInfo.AppName, dbInfo.UserId), metav1.DeleteOptions{})
}
