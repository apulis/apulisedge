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
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"time"
)

var logger = loggers.LogInstance()

func CreateEdgeNode(userInfo proto.ApulisHeader, req *nodemodule.CreateEdgeNodeReq) (*nodeentity.NodeBasicInfo, error) {
	node := &nodeentity.NodeBasicInfo{
		ClusterId:        userInfo.ClusterId,
		GroupId:          userInfo.GroupId,
		UserId:           userInfo.UserId,
		NodeName:         req.Name,
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

func ListEdgeNodes(userInfo proto.ApulisHeader, req *nodemodule.ListEdgeNodesReq) (*[]nodeentity.NodeBasicInfo, int, error) {
	var nodeInfos []nodeentity.NodeBasicInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Find(&nodeInfos)

	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	total = int(res.RowsAffected)
	return &nodeInfos, total, nil
}

func DescribeEdgeNode(userInfo proto.ApulisHeader, req *nodemodule.DescribeEdgeNodesReq) (*nodeentity.NodeBasicInfo, error) {
	var nodeInfo nodeentity.NodeBasicInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
		First(&nodeInfo)

	if res.Error != nil {
		return &nodeInfo, res.Error
	}

	return &nodeInfo, nil
}

// delete edge node
func DeleteEdgeNode(userInfo proto.ApulisHeader, req *nodemodule.DeleteEdgeNodeReq) error {
	var nodeInfo nodeentity.NodeBasicInfo

	// first: check if any deploy exist
	var total int64
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
		Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: get node and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
		First(&nodeInfo)
	if res.Error != nil {
		return res.Error
	}

	return nodeentity.DeleteNode(&nodeInfo)
}
