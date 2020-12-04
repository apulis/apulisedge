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
	var k8sInfo *nodeentity.NodeBasicInfo
	var dbInfo *nodeentity.NodeBasicInfo
	var err error

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

				k8sInfo, err = GetK8sNodeInfo(&nodeInfos[i])
				dbInfo = new(nodeentity.NodeBasicInfo)
				if err != nil && nodeInfos[i].Status != constants.StatusNotInstalled { // uninstall node
					logger.Infof("NodeTicker get k8sNodeInfo failed, so kick it. err = %v", err)
					dbInfo.ID = nodeInfos[i].ID
					dbInfo.UserId = nodeInfos[i].UserId
					dbInfo.UserName = nodeInfos[i].UserName
					dbInfo.NodeName = nodeInfos[i].NodeName
					dbInfo.Status = constants.StatusNotInstalled
					dbInfo.CreateAt = nodeInfos[i].CreateAt
					dbInfo.UpdateAt = time.Now()
					err = nodeentity.UpdateNode(dbInfo)
					if err != nil {
						logger.Infof("NodeTicker kick node failed, node = %s, err = %v", dbInfo.NodeName, err)
					} else {
						logger.Infof("NodeTicker kick node succ, node = %s", dbInfo.NodeName)
					}
				} else if err == nil && nodeInfos[i].Status == constants.StatusNotInstalled { // install node
					err = nodeentity.UpdateNode(k8sInfo)
					if err != nil {
						logger.Infof("NodeTicker install node failed, node = %s, err = %v", k8sInfo.NodeName, err)
					} else {
						logger.Infof("NodeTicker install node succ, node = %s", k8sInfo.NodeName)
					}
				}
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

func GetK8sNodeInfo(dbInfo *nodeentity.NodeBasicInfo) (*nodeentity.NodeBasicInfo, error) {
	var interIp string
	var outerIp string
	var nodeStatus string
	var roles string

	nodeInfo, err := utils.DescribeNode(dbInfo.NodeName)
	if err != nil {
		logger.Debugf("GetK8sNodeInfo DescribeNode failed. name = %s, err = %v", dbInfo.NodeName, err)
		return nil, err
	}

	// get address
	for _, addr := range nodeInfo.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			interIp = addr.Address
		} else if addr.Type == v1.NodeExternalIP {
			outerIp = addr.Address
		}
	}

	// get condition
	for _, cond := range nodeInfo.Status.Conditions {
		if cond.Type == v1.NodeReady {
			if cond.Status == v1.ConditionTrue {
				nodeStatus = constants.StatusOnline
			} else if cond.Status == v1.ConditionFalse {
				nodeStatus = constants.StatusOffline
			}
		}
	}

	// get roles
	for k, _ := range nodeInfo.Labels {
		if k == constants.AgentRoleKey {
			roles = strings.Join([]string{roles, constants.AgentRole}, ",")
		} else if k == constants.EdgeRoleKey {
			roles = strings.Join([]string{roles, constants.EdgeRole}, ",")
		}
	}
	roles = strings.TrimPrefix(roles, ",")
	roles = strings.TrimSuffix(roles, ",")

	k8sInfo := &nodeentity.NodeBasicInfo{
		ID:               dbInfo.ID,
		UserId:           dbInfo.UserId,
		UserName:         dbInfo.UserName,
		NodeName:         nodeInfo.Name,
		Status:           nodeStatus,
		Roles:            roles,
		ContainerRuntime: nodeInfo.Status.NodeInfo.ContainerRuntimeVersion,
		OsImage:          nodeInfo.Status.NodeInfo.OSImage,
		ProviderId:       nodeInfo.Spec.ProviderID,
		InterIp:          interIp,
		OuterIp:          outerIp,
		CreateAt:         dbInfo.CreateAt,
		UpdateAt:         time.Now(),
	}

	return k8sInfo, nil
}
