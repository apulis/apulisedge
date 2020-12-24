// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"context"
	"encoding/json"
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	applicationentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"time"
)

type statusHandler func(appDeployInfo *applicationentity.ApplicationDeployInfo)

// status transfer
var statusHandlerMap = map[string]statusHandler{
	constants.StatusInit:       handleStatusInit,
	constants.StatusDeploying:  handleStatusDeploying,
	constants.StatusRunning:    handleStatusRunning,
	constants.StatusAbnormal:   handleStatusAbnormal,
	constants.StatusDeleting:   handleStatusDeleting,
	constants.StatusUpdateInit: handleStatusUpdateInit,
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
				if _, ok := statusHandlerMap[appDeployInfos[i].Status]; ok {
					statusHandlerMap[appDeployInfos[i].Status](&appDeployInfos[i])
				} else {
					logger.Errorf("ApplicationTicker: No valid handler, status = %s", appDeployInfos[i].Status)
				}
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

func handleStatusInit(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var deployExist bool

	// check deploy status, if ok, transfer to StatusDeploying; if failed, retry
	_, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusInit GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
		// deploy to k8s
		err := CreateK8sDeployment(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusInit create deploy failed! err = %v", err)
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

	logger.Infof("handleStatusInit create deploy successful! status to %s", constants.StatusDeploying)
}

func handleStatusDeploying(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var deployExist bool

	// check deploy status
	deploy, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusDeploying GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
		appDeployInfo.Status = constants.StatusInit
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeploying update deployment failed!")
		}

		logger.Infof("handleStatusDeploying deployment disappeared! return to init status!")
		return
	}

	if deploy.Status.ReadyReplicas == deploy.Status.Replicas {
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
	var deployExist bool

	// check deploy status
	deploy, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusAbnormal GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
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

	if deploy.Status.ReadyReplicas == deploy.Status.Replicas {
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
	var deployExist bool

	// check deploy status
	deploy, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusRunning GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
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

	if deploy.Status.ReadyReplicas != deploy.Status.Replicas {
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

func handleStatusUpdateInit(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var deployExist bool

	// check deploy status
	_, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusUpdateInit GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
		appDeployInfo.Status = constants.StatusInit
		appDeployInfo.UpdateAt = time.Now()
		err := applicationentity.UpdateAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusUpdateInit update deployment failed!")
		}

		logger.Infof("handleStatusUpdateInit deployment disappeared! return to init status!")
		return
	}

	err = UpdateK8sDeployment(appDeployInfo)
	if err != nil {
		logger.Infof("handleStatusUpdateInit update deploy failed! err = %v", err)
		return
	}

	// update status
	appDeployInfo.Status = constants.StatusDeploying
	appDeployInfo.UpdateAt = time.Now()
	err = applicationentity.UpdateAppDeploy(appDeployInfo)
	if err != nil {
		logger.Infof("handleStatusUpdateInit update deployment failed when deploying!")
		return
	}

	logger.Infof("handleStatusUpdateInit update deploy successful! status to %s", constants.StatusDeploying)
}

func handleStatusDeleting(appDeployInfo *applicationentity.ApplicationDeployInfo) {
	var deployExist bool

	// check deploy status
	_, err := GetK8sDeployment(appDeployInfo)
	if err == nil {
		deployExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			deployExist = false
		} else {
			logger.Infof("handleStatusRunning GetK8sDeployment! err = %v", err)
			return
		}
	}

	if !deployExist {
		err = applicationentity.DeleteAppDeploy(appDeployInfo)
		if err != nil {
			logger.Infof("handleStatusDeleting, DeleteAppDeploy err, err = %v", err)
			return
		}
		logger.Infof("handleStatusDeleting, DeleteAppDeploy succ")
		return
	}

	err = DeleteK8sDeployment(appDeployInfo)
	if err != nil { // if failed, try next time
		logger.Infof("handleStatusDeleting, DeleteK8sDeployment err, err = %v", err)
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

func CreateK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) error {
	var appVerInfo applicationentity.ApplicationVersionInfo

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		return err
	}

	deployClient, err := clu.GetDeploymentClient(constants.DefaultNamespace)
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

	var appNetwork appmodule.CreateAppNetwork
	err = json.Unmarshal([]byte(appVerInfo.Network), &appNetwork)
	if err != nil {
		return err
	}
	logger.Infof("appNetwork = %+v", appNetwork)

	// host network
	hn := false
	var ports []corev1.ContainerPort
	if appNetwork.Type == appmodule.NetworkTypeHost {
		hn = true
	}
	if appNetwork.Type == appmodule.NetworkTypePortMapping {
		for _, v := range appNetwork.PortMappings {
			ports = append(ports, corev1.ContainerPort{
				ContainerPort: int32(v.ContainerPort),
				HostPort:      int32(v.HostPort),
			})
		}
	}

	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbInfo.DeployUUID,
		},
		Spec: v1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1), // 指定副本数
			Selector: &metav1.LabelSelector{ // 指定标签
				MatchLabels: map[string]string{
					"app": dbInfo.AppName,
				},
			},
			Template: corev1.PodTemplateSpec{ // 容器模板
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": dbInfo.AppName,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": dbInfo.UniqueName,
					},
					HostNetwork:   hn,
					RestartPolicy: corev1.RestartPolicy(appVerInfo.RestartPolicy),
					Containers: []corev1.Container{
						{
							Name:  dbInfo.ContainerUUID,
							Image: appVerInfo.ContainerImagePath,
							Ports: ports,
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

func UpdateK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) error {
	var appVerInfo applicationentity.ApplicationVersionInfo

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		return err
	}

	deployClient, err := clu.GetDeploymentClient(constants.DefaultNamespace)
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

	deployment, err := deployClient.Get(context.Background(), dbInfo.DeployUUID, metav1.GetOptions{})
	if err != nil {
		logger.Infof("err to get deployment, app = %s, version = %s, node = %s, deployId = %s",
			dbInfo.AppName, dbInfo.Version, dbInfo.NodeName, dbInfo.DeployUUID)
		return err
	}

	// add rolling update strategy
	maxSur := intstr.FromInt(1)
	maxUna := intstr.FromInt(1)
	for i, _ := range deployment.Spec.Template.Spec.Containers {
		if deployment.Spec.Template.Spec.Containers[i].Name == dbInfo.ContainerUUID {
			deployment.Spec.Template.Spec.Containers[i].Image = appVerInfo.ContainerImagePath
		}
	}

	deployment.Spec.MinReadySeconds = 5
	deployment.Spec.Strategy = v1.DeploymentStrategy{
		Type: v1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &v1.RollingUpdateDeployment{
			MaxSurge:       &maxSur,
			MaxUnavailable: &maxUna,
		},
	}

	_, err = deployClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) (*v1.Deployment, error) {
	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("GetK8sDeployment, can`t find cluster %d", dbInfo.ClusterId)
		return nil, err
	}

	deployClient, err := clu.GetDeploymentClient(constants.DefaultNamespace)
	if err != nil {
		return nil, err
	}

	deploy, err := deployClient.Get(context.Background(), dbInfo.DeployUUID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

func DeleteK8sDeployment(dbInfo *applicationentity.ApplicationDeployInfo) error {
	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("DeleteK8sDeployment, can`t find cluster %d", dbInfo.ClusterId)
		return err
	}

	deployClient, err := clu.GetDeploymentClient(constants.DefaultNamespace)
	if err != nil {
		return err
	}

	return deployClient.Delete(context.Background(), dbInfo.DeployUUID, metav1.DeleteOptions{})
}
