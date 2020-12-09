// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type statusHandler func(appDeployInfo *nodeentity.NodeBasicInfo)

// status transfer
var statusHandlerMap = map[string]statusHandler{
	constants.StatusNotInstalled: handleStatusNotInstalled,
	constants.StatusOnline:       handleStatusOnline,
	constants.StatusOffline:      handleStatusOffline,
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
		} else {
			for i := 0; i < int(res.RowsAffected); i++ {
				logger.Debugf("NodeTicker handle node = %v", nodeInfos[i])
				statusHandlerMap[nodeInfos[i].Status](&nodeInfos[i])
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

// keep uninstalled or to online/offline
func handleStatusNotInstalled(dbInfo *nodeentity.NodeBasicInfo) {
	var interIp string
	var outerIp string
	var nodeStatus string
	var roles string

	// first: get node info from k8s
	nodeK8sInfo, err := utils.DescribeNode(dbInfo.NodeName)
	if err != nil {
		logger.Debugf("handleStatusNotInstalled DescribeNode failed. name = %s, err = %v", dbInfo.NodeName, err)
		return
	}

	// second: create new node and install node
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
				nodeStatus = constants.StatusOnline
			} else if cond.Status == v1.ConditionFalse {
				nodeStatus = constants.StatusOffline
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

	newDbInfo := &nodeentity.NodeBasicInfo{
		ID:               dbInfo.ID,
		ClusterId:        dbInfo.ClusterId,
		GroupId:          dbInfo.GroupId,
		UserId:           dbInfo.UserId,
		NodeName:         nodeK8sInfo.Name,
		Status:           nodeStatus,
		Roles:            roles,
		ContainerRuntime: nodeK8sInfo.Status.NodeInfo.ContainerRuntimeVersion,
		OsImage:          nodeK8sInfo.Status.NodeInfo.OSImage,
		InterIp:          interIp,
		OuterIp:          outerIp,
		CreateAt:         dbInfo.CreateAt,
		UpdateAt:         time.Now(),
	}

	err = nodeentity.UpdateNode(newDbInfo)
	if err != nil {
		logger.Infof("NodeTicker install node failed, node = %s, err = %v", dbInfo.NodeName, err)
	} else {
		logger.Infof("NodeTicker install node succ, node = %s", dbInfo.NodeName)
	}
}

// keep online or to offline/uninstalled
func handleStatusOnline(dbInfo *nodeentity.NodeBasicInfo) {
	var newDbInfo nodeentity.NodeBasicInfo
	var nodeStatus string

	nodeK8sInfo, err := utils.DescribeNode(dbInfo.NodeName)
	// if err get k8s info, turn to uninstalled
	if err != nil {
		logger.Debugf("handleStatusOnline DescribeNode failed. name = %s, err = %v", dbInfo.NodeName, err)
		// turn to uninstalled
		newDbInfo = *dbInfo
		newDbInfo.Status = constants.StatusNotInstalled
		newDbInfo.UpdateAt = time.Now()
		err = nodeentity.UpdateNode(&newDbInfo)
		if err != nil {
			logger.Infof("handleStatusOnline UpdateNode failed, node = %s, err = %v", dbInfo.NodeName, err)
		} else {
			logger.Infof("handleStatusOnline UpdateNode succ, node = %s, status = %s", dbInfo.NodeName, newDbInfo.Status)
		}
	} else { // check if not ready
		for _, cond := range nodeK8sInfo.Status.Conditions {
			if cond.Type == v1.NodeReady {
				if cond.Status == v1.ConditionTrue {
					nodeStatus = constants.StatusOnline
				} else if cond.Status == v1.ConditionFalse {
					nodeStatus = constants.StatusOffline
				}
			}
		}

		if nodeStatus == constants.StatusOffline {
			newDbInfo = *dbInfo
			newDbInfo.Status = constants.StatusOffline
			newDbInfo.UpdateAt = time.Now()
			err = nodeentity.UpdateNode(&newDbInfo)
			if err != nil {
				logger.Infof("handleStatusOnline UpdateNode failed, node = %s, err = %v", dbInfo.NodeName, err)
			} else {
				logger.Infof("handleStatusOnline UpdateNode succ, node = %s, status = %s", dbInfo.NodeName, newDbInfo.Status)
			}
		}
	}
}

func handleStatusOffline(dbInfo *nodeentity.NodeBasicInfo) {

}
