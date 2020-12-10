// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
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
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ?", req.ClusterId, req.GroupId, req.UserId).
		Find(&nodeInfos)

	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	return &nodeInfos, int(res.RowsAffected), nil
}

func DescribeEdgeNode(req *nodemodule.DescribeEdgeNodesReq) (*nodeentity.NodeBasicInfo, error) {
	var nodeInfo nodeentity.NodeBasicInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", req.ClusterId, req.GroupId, req.UserId, req.NodeName).
		First(&nodeInfo)

	if res.Error != nil {
		return &nodeInfo, res.Error
	}

	return &nodeInfo, nil
}

// delete edge node
func DeleteEdgeNode(req *nodemodule.DeleteEdgeNodeReq) error {
	var nodeInfo nodeentity.NodeBasicInfo

	// first: check if any deploy exist
	var total int64
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", req.ClusterId, req.GroupId, req.UserId, req.NodeName).
		Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: get node and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", req.ClusterId, req.GroupId, req.UserId, req.NodeName).
		First(&nodeInfo)
	if res.Error != nil {
		return res.Error
	}

	return nodeentity.DeleteNode(&nodeInfo)
}
