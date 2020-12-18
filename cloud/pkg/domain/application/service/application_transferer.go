// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	applicationentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	"github.com/satori/go.uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type statusHandler func(appDeployInfo *applicationentity.ApplicationDeployInfo)

// status transfer
var statusHandlerMap = map[string]statusHandler{
	constants.StatusInit:      handleStatusInit,
	constants.StatusDeploying: handleStatusDeploying,
	constants.StatusRunning:   handleStatusRunning,
	constants.StatusAbnormal:  handleStatusAbnormal,
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
			continue
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
	var podExist bool

	// check pod status, if ok, transfer to StatusDeploying; if failed, retry
	_, err := GetK8sPod(appDeployInfo)
	if err == nil {
		podExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			podExist = false
		} else {
			logger.Infof("handleStatusInit GetK8sPod! err = %v", err)
			return
		}
	}

	if !podExist {
		// deploy to k8s
		err := CreateK8sPod(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusInit create pod failed! err = %v", err)
			return
		}
	}

	// update status
	appDeployInfo.Status = constants.StatusDeploying
	appDeployInfo.UpdateAt = time.Now()
	err = applicationentity.UpdateAppDeploy(appDeployInfo)
	if err != nil {
		logger.Infof("handleStatusInit update deployment failed when deploying!")
		return
	}

	logger.Infof("handleStatusInit create pod succ! status to %s", constants.StatusDeploying)
}

func handleStatusDeploying(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var podExist bool

	// check pod status
	pod, err := GetK8sPod(appDeployInfo)
	if err == nil {
		podExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			podExist = false
		} else {
			logger.Infof("handleStatusDeploying GetK8sPod! err = %v", err)
			return
		}
	}

	if !podExist {
		appDeployInfo.Status = constants.StatusInit
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
		}

		logger.Infof("handleStatusDeploying deployment disappeared! return to init status!")
		return
	}

	if pod.Status.Phase == corev1.PodRunning {
		appDeployInfo.Status = constants.StatusRunning
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
			return
		}
		logger.Infof("handleStatusDeploying UpdateAppDeploy succ! status to %s", appDeployInfo.Status)
	} else {
		appDeployInfo.Status = constants.StatusAbnormal
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
			return
		}
		logger.Infof("handleStatusDeploying UpdateAppDeploy succ! status to %s", appDeployInfo.Status)
	}
}

func handleStatusAbnormal(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var podExist bool

	// check pod status
	pod, err := GetK8sPod(appDeployInfo)
	if err == nil {
		podExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			podExist = false
		} else {
			logger.Infof("handleStatusAbnormal GetK8sPod! err = %v", err)
			return
		}
	}

	if !podExist {
		appDeployInfo.Status = constants.StatusInit
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusAbnormal update deployment failed!")
			return
		}

		logger.Infof("handleStatusAbnormal deployment disappeared! return to init status!")
		return
	}

	if pod.Status.Phase == corev1.PodRunning {
		appDeployInfo.Status = constants.StatusRunning
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusAbnormal update deployment failed!")
			return
		}

		logger.Infof("handleStatusDeploying UpdateAppDeploy succ! status to %s", appDeployInfo.Status)
	}
}

func handleStatusRunning(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var podExist bool

	// check pod status
	pod, err := GetK8sPod(appDeployInfo)
	if err == nil {
		podExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			podExist = false
		} else {
			logger.Infof("handleStatusRunning GetK8sPod! err = %v", err)
			return
		}
	}

	if !podExist {
		appDeployInfo.Status = constants.StatusInit
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusRunning update deployment failed!")
			return
		}

		logger.Infof("handleStatusRunning deployment disappeared! return to init status!")
		return
	}

	if pod.Status.Phase != corev1.PodRunning {
		appDeployInfo.Status = constants.StatusAbnormal
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
			return
		}

		logger.Infof("handleStatusDeploying UpdateAppDeploy succ! status to %s", appDeployInfo.Status)
	}
}

func handleStatusDeleting(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var podExist bool

	// check pod status
	_, err := GetK8sPod(appDeployInfo)
	if err == nil {
		podExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			podExist = false
		} else {
			logger.Infof("handleStatusRunning GetK8sPod! err = %v", err)
			return
		}
	}

	if !podExist {
		err = applicationentity.DeleteAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeleting, DeleteAppDeploy err, err = %v", err)
			return
		}
		logger.Infof("handleStatusDeleting, DeleteAppDeploy succ")
		return
	}

	err = DeleteK8sPod(appDeployInfo)
	if err != nil { // if failed, try next time
		logger.Infof("handleStatusDeleting, DeleteK8sPod err, err = %v", err)
		return
	}

	// db record delete
	err = applicationentity.DeleteAppDeploy(appDeployInfo)
	if err != nil {
		logger.Infof("handleStatusDeleting, DeleteAppDeploy err, err = %v", err)
		return
	}
	logger.Infof("handleStatusDeleting, DeleteAppDeploy succ")
}

func CreateK8sPod(dbInfo *applicationentity.ApplicationDeployInfo) error {
	var appVerInfo applicationentity.ApplicationVersionInfo

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		return err
	}

	podClient, err := clu.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return err
	}

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			dbInfo.ClusterId, dbInfo.GroupId, dbInfo.UserId, dbInfo.AppName, dbInfo.Version).
		First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbInfo.DeployUUID,
		},
		Spec: corev1.PodSpec{
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": dbInfo.NodeName,
			},
			Containers: []corev1.Container{
				{
					Name:  uuid.NewV4().String(),
					Image: appVerInfo.ContainerImagePath,
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
	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("GetK8sPod, can`t find cluster %d", dbInfo.ClusterId)
		return nil, err
	}

	podClient, err := clu.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return nil, err
	}

	pod, err := podClient.Get(context.Background(), dbInfo.DeployUUID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func DeleteK8sPod(dbInfo *applicationentity.ApplicationDeployInfo) error {
	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("DeleteK8sPod, can`t find cluster %d", dbInfo.ClusterId)
		return err
	}

	podClient, err := clu.GetPodClient(constants.DefaultNamespace)
	if err != nil {
		return err
	}

	return podClient.Delete(context.Background(), dbInfo.DeployUUID, metav1.DeleteOptions{})
}
