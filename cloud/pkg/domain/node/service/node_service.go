// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"fmt"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"time"
)

var logger = loggers.LogInstance()

func CreateEdgeNode(req *nodemodule.CreateEdgeNodeReq) (*nodeentity.NodeBasicInfo, error) {
	node := &nodeentity.NodeBasicInfo{
		ClusterId:        req.ClusterId,
		GroupId:          req.GroupId,
		UserId:           req.UserId,
		NodeName:         req.NodeName,
		Status:           constants.StatusNotInstalled,
		Roles:            "",
		ContainerRuntime: "",
		OsImage:          "",
		InterIp:          "",
		OuterIp:          "",
		CreateAt:         time.Now(),
		UpdateAt:         time.Now(),
	}

	return node, nodeentity.CreateNode(node)
}

func ListEdgeNodes(req *nodemodule.ListEdgeNodesReq) (*[]nodeentity.NodeBasicInfo, int, error) {
	var nodeInfos []nodeentity.NodeBasicInfo

	total := 0
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' ", req.ClusterId, req.GroupId, req.UserId)
	res := apulisdb.Db.Offset(req.PageNum).Limit(req.PageSize).Where(whereQueryStr).Find(&nodeInfos)

	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	return &nodeInfos, int(res.RowsAffected), nil
}

func DescribeEdgeNode(req *nodemodule.DescribeEdgeNodesReq) (*nodeentity.NodeBasicInfo, error) {
	var nodeInfo nodeentity.NodeBasicInfo

	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and NodeName = '%s'", req.ClusterId, req.GroupId, req.UserId, req.NodeName)
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
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and NodeName = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.NodeName)
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).Where(whereQueryStr).Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: check if any node exist
	whereQueryStr = fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and NodeName = '%s'", req.ClusterId, req.GroupId, req.UserId, req.NodeName)
	apulisdb.Db.Model(&nodeentity.NodeBasicInfo{}).Where(whereQueryStr).Count(&total)
	if total != 0 {
		return nodemodule.ErrNodeNotExist
	}

	nodeInfo, err := nodeentity.GetNode(req.ClusterId, req.GroupId, req.UserId, req.NodeName)
	if err != nil {
		return err
	}

	return nodeentity.DeleteNode(nodeInfo)
}
