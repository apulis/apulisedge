// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"github.com/apulis/ApulisEdge/loggers"
	constants "github.com/apulis/ApulisEdge/node"
	nodeentity "github.com/apulis/ApulisEdge/node/entity"
	"github.com/apulis/ApulisEdge/utils"
	v1 "k8s.io/api/core/v1"
	"time"
)

var logger = loggers.Log

func ListEdgeNodes() ([]*nodeentity.NodeBasicInfo, error) {
	var interIp string
	var outerIp string
	var nodeStatus string

	nodeList, err := utils.ListNodes()
	if err != nil {
		logger.Info("Failed to listEdgeNodes, err = [%v]", err)
		return nil, err
	}

	nodes := make([]*nodeentity.NodeBasicInfo, 0)
	for _, v := range nodeList.Items {
		if !isEdgeNode(&v) {
			continue
		}

		// get address
		for _, addr := range v.Status.Addresses {
			if addr.Type == v1.NodeInternalIP {
				interIp = addr.Address
			} else if addr.Type == v1.NodeExternalIP {
				outerIp = addr.Address
			}
		}

		// get condition
		for _, cond := range v.Status.Conditions {
			if cond.Type == v1.NodeReady {
				if cond.Status == v1.ConditionTrue {
					nodeStatus = constants.StatusOnline
				} else if cond.Status == v1.ConditionFalse {
					nodeStatus = constants.StatusOffline
				}
			}
		}

		nodes = append(nodes, &nodeentity.NodeBasicInfo{
			Name:             v.Name,
			Status:           nodeStatus,
			Roles:            "",
			ContainerRuntime: v.Status.NodeInfo.ContainerRuntimeVersion,
			OsImage:          v.Status.NodeInfo.OSImage,
			ProviderId:       v.Spec.ProviderID,
			InterIp:          interIp,
			OuterIp:          outerIp,
			CreateTime:       time.Now().String(),
		})

	}

	return nodes, nil
}

func isEdgeNode(v *v1.Node) bool {
	if _, ok := v.Labels[constants.EdgeRoleKey]; ok {
		return true
	}

	return false
}

func DescribeEdgeNode(name string) (*nodeentity.NodeDetailInfo, error) {
	var interIp string
	var outerIp string
	var nodeStatus string

	nodeInfo, err := utils.DescribeNode(name)
	if err != nil {
		logger.Info("Failed to listEdgeNodes, err = [%v]", err)
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

	return &nodeentity.NodeDetailInfo{
		Name:             nodeInfo.Name,
		Status:           nodeStatus,
		Roles:            "",
		ContainerRuntime: nodeInfo.Status.NodeInfo.ContainerRuntimeVersion,
		OsImage:          nodeInfo.Status.NodeInfo.OSImage,
		ProviderId:       nodeInfo.Spec.ProviderID,
		InterIp:          interIp,
		OuterIp:          outerIp,
		CreateTime:       time.Now().String(),
	}, nil
}
