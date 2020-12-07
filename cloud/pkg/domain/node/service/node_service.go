// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"fmt"
	"time"

	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
)

var logger = loggers.LogInstance()

func CreateEdgeNode(userId int64, userName string, nodeName string) (*nodeentity.NodeBasicInfo, error) {
	node := &nodeentity.NodeBasicInfo{
		UserId:           userId,
		UserName:         userName,
		NodeName:         nodeName,
		Status:           constants.StatusNotInstalled,
		Roles:            "",
		ContainerRuntime: "",
		OsImage:          "",
		ProviderId:       "",
		InterIp:          "",
		OuterIp:          "",
		CreateAt:         time.Now(),
		UpdateAt:         time.Now(),
	}

	return node, nodeentity.CreateNode(node)
}

func ListEdgeNodes(userId int64, offset int, limit int) (*[]nodeentity.NodeBasicInfo, int, error) {
	var nodeInfos []nodeentity.NodeBasicInfo

	total := 0
	whereQueryStr := fmt.Sprintf("UserId = '%s' ", userId)
	res := apulisdb.Db.Offset(offset).Limit(limit).Where(whereQueryStr).Find(&nodeInfos)

	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	return &nodeInfos, int(res.RowsAffected), nil
}

func DescribeEdgeNode(userId int64, nodeName string) (*nodeentity.NodeBasicInfo, error) {
	var nodeInfo nodeentity.NodeBasicInfo

	whereQueryStr := fmt.Sprintf("UserId = '%s' and NodeName = '%s'", userId, nodeName)
	res := apulisdb.Db.Where(whereQueryStr).First(&nodeInfo)

	if res.Error != nil {
		return &nodeInfo, res.Error
	}

	return &nodeInfo, nil
}

// delete edge node
func DeleteEdgeNode(req *nodemodule.DeleteEdgeNodeReq) error {
	// first: check if any deploy exist
	var total int64
	whereQueryStr := fmt.Sprintf("UserId = '%s' and NodeName = '%s'", req.UserId, req.NodeName)
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).Where(whereQueryStr).Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: check if any node exist
	whereQueryStr = fmt.Sprintf("UserId = '%s' and NodeName = '%s'", req.UserId, req.NodeName)
	apulisdb.Db.Model(&nodeentity.NodeBasicInfo{}).Where(whereQueryStr).Count(&total)
	if total == 0 {
		return nodemodule.ErrNodeNotExist
	}

	nodeInfo, err := nodeentity.GetNode(req.UserId, req.NodeName)
	if err != nil {
		return err
	}

	return nodeentity.DeleteNode(nodeInfo)
}
