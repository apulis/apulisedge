// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeticker

import (
	"context"
	"strings"
	"time"

	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logger = loggers.LogInstance()

type statusHandler func(appDeployInfo *nodeentity.NodeBasicInfo)

// status transfer
var statusHandlerMap = map[string]statusHandler{
	constants.StatusNotInstalled: handleStatusNotInstalled,
	constants.StatusOnline:       handleStatusOnline,
	constants.StatusOffline:      handleStatusOffline,
	constants.StatusDeleting:     handleStatusDelete,
	constants.StatusInstalling:   handleStatusInstalling,
}

// CreateNodeCheckLoop transferer of edge node status
func CreateNodeTickerLoop(ctx context.Context, config *configs.EdgeCloudConfig) {
	duration := time.Duration(config.Portal.NodeCheckerInterval) * time.Second
	checkTicker := time.NewTimer(duration)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("CreateNodeTickerLoop was terminated")
			return
		case <-checkTicker.C:
			NodeTicker(config)
			checkTicker.Reset(duration)
		}
	}
}

func NodeTicker(config *configs.EdgeCloudConfig) {
	var nodeInfos []nodeentity.NodeBasicInfo
	var totalTmp int64
	var total int
	offset := 0

	apulisdb.Db.Model(&nodeentity.NodeBasicInfo{}).Count(&totalTmp)
	total = int(totalTmp)

	logger.Debugf("NodeTicker total node count = %d", total)
	if total == 0 {
		return
	} else if total < constants.TransferCountEach {
		total = constants.TransferCountEach
	}

	for total >= constants.TransferCountEach {
		res := apulisdb.Db.Offset(offset).Limit(constants.TransferCountEach).Find(&nodeInfos)
		if res.Error != nil {
			logger.Errorf("query node failed. err = %v", res.Error)
			continue
		} else {
			for i := 0; i < int(res.RowsAffected); i++ {
				logger.Debugf("NodeTicker handle node = %v", nodeInfos[i])
				if _, ok := statusHandlerMap[nodeInfos[i].Status]; ok {
					statusHandlerMap[nodeInfos[i].Status](&nodeInfos[i])
				} else {
					logger.Errorf("NodeTicker: No valid handler, status = %s", nodeInfos[i].Status)
				}
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

// keep uninstalled or to online/offline
func handleStatusNotInstalled(dbInfo *nodeentity.NodeBasicInfo) {
	var newDbInfo nodeentity.NodeBasicInfo

	var interIp string
	var outerIp string
	var curNodeStatus string
	var roles string

	var nodeExist bool

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("handleStatusNotInstalled, can`t find cluster %d", dbInfo.ClusterId)
		return
	}

	// get node info from k8s
	nodeK8sInfo, err := clu.DescribeNode(dbInfo.UniqueName)
	if err == nil {
		nodeExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			nodeExist = false
		} else {
			logger.Infof("handleStatusNotInstalled DescribeNode failed! nodeName = %s, uniName = %s, err = %v",
				dbInfo.NodeName, dbInfo.UniqueName, err)
			return
		}
	}

	// node not exist in k8s, try next time
	if !nodeExist {
		logger.Debugf("handleStatusNotInstalled name %s not exist in kubernetes. nodeName = %s", dbInfo.UniqueName, dbInfo.NodeName)
		return
	}

	var cpuCore int
	var mem int64
	if _, ok := nodeK8sInfo.Status.Capacity[v1.ResourceCPU]; ok {
		q := nodeK8sInfo.Status.Capacity[v1.ResourceCPU]
		v, b := (&q).AsInt64()
		if b {
			cpuCore = int(v)
		}
	}

	if _, ok := nodeK8sInfo.Status.Capacity[v1.ResourceMemory]; ok {
		q := nodeK8sInfo.Status.Capacity[v1.ResourceMemory]
		v, b := (&q).AsInt64()
		if b {
			mem = v
		}
	}

	// create new node and install node
	// get address
	for _, addr := range nodeK8sInfo.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			interIp = addr.Address
		} else if addr.Type == v1.NodeExternalIP {
			outerIp = addr.Address
		}
	}

	// get condition
	for _, cond := range nodeK8sInfo.Status.Conditions {
		if cond.Type == v1.NodeReady {
			if cond.Status == v1.ConditionTrue {
				curNodeStatus = constants.StatusOnline
			} else if cond.Status == v1.ConditionFalse || cond.Status == v1.ConditionUnknown {
				curNodeStatus = constants.StatusOffline
			}
		}
	}

	// get roles
	for k, _ := range nodeK8sInfo.Labels {
		if k == constants.AgentRoleKey {
			roles = strings.Join([]string{roles, constants.AgentRole}, ",")
		} else if k == constants.EdgeRoleKey {
			roles = strings.Join([]string{roles, constants.EdgeRole}, ",")
		}
	}
	roles = strings.TrimPrefix(roles, ",")
	roles = strings.TrimSuffix(roles, ",")

	newDbInfo = *dbInfo
	newDbInfo.Status = curNodeStatus
	newDbInfo.Roles = roles
	newDbInfo.ContainerRuntime = nodeK8sInfo.Status.NodeInfo.ContainerRuntimeVersion
	newDbInfo.OsImage = nodeK8sInfo.Status.NodeInfo.OSImage
	newDbInfo.InterIp = interIp
	newDbInfo.OuterIp = outerIp
	newDbInfo.CpuCore = cpuCore
	newDbInfo.Memory = mem
	newDbInfo.Arch = nodeK8sInfo.Status.NodeInfo.Architecture
	newDbInfo.UpdateAt = time.Now()

	// label一下，方便查看
	err = clu.LabelNode(nodeK8sInfo, newDbInfo.ClusterId, newDbInfo.GroupId, newDbInfo.UserId, newDbInfo.NodeName)
	if err != nil {
		logger.Infof("NodeTicker label node failed, node = %s, uniq = %s, err = %v", dbInfo.NodeName, dbInfo.UniqueName, err)
		return
	}

	err = nodeentity.UpdateNode(&newDbInfo)
	if err != nil {
		logger.Infof("NodeTicker install node failed, node = %s, uniq = %s, err = %v", dbInfo.NodeName, dbInfo.UniqueName, err)
	} else {
		logger.Infof("NodeTicker install node succ, node = %s, uniq = %s, turn to status %s", dbInfo.NodeName, dbInfo.UniqueName, newDbInfo.Status)
	}
}

// keep online or to offline/uninstalled
func handleStatusOnline(dbInfo *nodeentity.NodeBasicInfo) {
	handleStatusInstalled(dbInfo)
}

func handleStatusOffline(dbInfo *nodeentity.NodeBasicInfo) {
	handleStatusInstalled(dbInfo)
}

func handleStatusDelete(dbInfo *nodeentity.NodeBasicInfo) {
	var nodeExist bool

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("handleStatusDelete, can`t find cluster %d", dbInfo.ClusterId)
		return
	}

	// get node info from k8s
	nodeK8sInfo, err := clu.DescribeNode(dbInfo.UniqueName)
	if err == nil {
		nodeExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			nodeExist = false
		} else {
			logger.Infof("handleStatusDelete DescribeNode failed! nodeName = %s, uniq = %s, err = %v",
				dbInfo.NodeName, dbInfo.UniqueName, err)
			return
		}
	}

	if !nodeExist {
		err = nodeentity.DeleteNode(dbInfo)
		if err != nil {
			logger.Infof("handleStatusDelete Delete From DB failed! nodeName = %s, uniq = %s, err = %v",
				dbInfo.NodeName, dbInfo.UniqueName, err)
		}
	} else {
		err = clu.DeleteNode(nodeK8sInfo)
		if err != nil {
			logger.Infof("handleStatusDelete Delete From k8s failed! nodeName = %s, uniq = %s, err = %v",
				dbInfo.NodeName, dbInfo.UniqueName, err)
		}
	}
}

func handleStatusInstalled(dbInfo *nodeentity.NodeBasicInfo) {
	var newDbInfo nodeentity.NodeBasicInfo
	var curNodeStatus string

	var nodeExist bool

	clu, err := cluster.GetCluster(dbInfo.ClusterId)
	if err != nil {
		logger.Infof("handleStatusInstalled, can`t find cluster %d", dbInfo.ClusterId)
		return
	}

	// get node info from k8s
	nodeK8sInfo, err := clu.DescribeNode(dbInfo.UniqueName)
	if err == nil {
		nodeExist = true
	} else {
		if errors.ReasonForError(err) == metav1.StatusReasonNotFound {
			nodeExist = false
		} else {
			logger.Infof("handleStatusInstalled DescribeNode failed! nodeName = %s, uniq = %s, err = %v",
				dbInfo.NodeName, dbInfo.UniqueName, err)
			// TODO try many times, if failed, turn to offline
			return
		}
	}

	// not exist, turn to init status
	if !nodeExist {
		newDbInfo = *dbInfo
		newDbInfo.Status = constants.StatusNotInstalled
		newDbInfo.UpdateAt = time.Now()
		err = nodeentity.UpdateNode(&newDbInfo)
		if err != nil {
			logger.Infof("handleStatusInstalled UpdateNode failed, node = %s, uniq = %s, err = %v", dbInfo.NodeName, dbInfo.UniqueName, err)
		} else {
			logger.Infof("handleStatusInstalled UpdateNode succ, node = %s, uniq = %s, turn to status = %s",
				dbInfo.NodeName, dbInfo.UniqueName, newDbInfo.Status)
		}
		return
	} else {
		for _, cond := range nodeK8sInfo.Status.Conditions {
			if cond.Type == v1.NodeReady {
				if cond.Status == v1.ConditionTrue {
					curNodeStatus = constants.StatusOnline
				} else {
					curNodeStatus = constants.StatusOffline
				}
			}
		}

		newDbInfo = *dbInfo
		newDbInfo.UpdateAt = time.Now()
		if dbInfo.Status == constants.StatusOnline && curNodeStatus == constants.StatusOffline {
			newDbInfo.Status = curNodeStatus
			err = nodeentity.UpdateNode(&newDbInfo)
			if err != nil {
				logger.Infof("handleStatusInstalled UpdateNode failed, node = %s, uniq = %s, err = %v", dbInfo.NodeName, dbInfo.UniqueName, err)
			} else {
				logger.Infof("handleStatusInstalled UpdateNode succ, node = %s, uniq = %s, turn to status = %s",
					dbInfo.NodeName, dbInfo.UniqueName, newDbInfo.Status)
			}
		} else if dbInfo.Status == constants.StatusOffline && curNodeStatus == constants.StatusOnline {
			newDbInfo.Status = curNodeStatus
			err = nodeentity.UpdateNode(&newDbInfo)
			if err != nil {
				logger.Infof("handleStatusInstalled UpdateNode failed, node = %s, uniq = %s, err = %v", dbInfo.NodeName, dbInfo.UniqueName, err)
			} else {
				logger.Infof("handleStatusInstalled UpdateNode succ, node = %s, uniq = %s, turn to status = %s",
					dbInfo.NodeName, dbInfo.UniqueName, newDbInfo.Status)
			}
		}
	}
}

func handleStatusInstalling(dbInfo *nodeentity.NodeBasicInfo) {
}
