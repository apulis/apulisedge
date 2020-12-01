// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"context"
	"fmt"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"strings"

	//"github.com/apulis/ApulisEdge/cloud/utils"
	//v1 "k8s.io/api/core/v1"
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
			logger.Infof("CreateNodeCheckLoop was terminated")
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
	var err error

	whereQueryStr := fmt.Sprintf("Status = '%s'", constants.StatusNotInstalled)
	apulisdb.Db.Model(&nodeentity.NodeBasicInfo{}).Where(whereQueryStr).Count(&totalTmp)
	total = int(totalTmp)
	if total == 0 {
		return
	} else if total < constants.TransferCountEach {
		total = constants.TransferCountEach
	}

	for total >= constants.TransferCountEach {
		res := apulisdb.Db.Offset(offset).Limit(constants.TransferCountEach).Where(whereQueryStr).Find(&nodeInfos)
		if res.Error != nil {
			logger.Errorf("query node failed. queryStr = %s, err = %v", whereQueryStr, res.Error)
		} else {
			if config.DebugModel {
				// debug print
				for i := 0; i < int(res.RowsAffected); i++ {
					logger.Infof("node ===> %s", nodeInfos[i].Name)
				}
			}
			for i := 0; i < int(res.RowsAffected); i++ {
				k8sInfo, err = NodeUpdate(&nodeInfos[i])
				if err != nil {
					logger.Infof("NodeTicker update node failed! err = %v", err)
				} else {
					logger.Infof("NodeTicker update node success! node = %s, status = %s", k8sInfo.Name, k8sInfo.Status)
				}
			}
		}

		offset += constants.TransferCountEach
		total -= constants.TransferCountEach
	}
}

func NodeUpdate(dbInfo *nodeentity.NodeBasicInfo) (*nodeentity.NodeBasicInfo, error) {
	var interIp string
	var outerIp string
	var nodeStatus string
	var roles string

	nodeInfo, err := utils.DescribeNode(dbInfo.Name)
	if err != nil {
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
		Name:             nodeInfo.Name,
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

	k8sInfo.ID = dbInfo.ID
	k8sInfo.UserId = dbInfo.UserId
	return k8sInfo, nodeentity.UpdateNode(k8sInfo)
}
