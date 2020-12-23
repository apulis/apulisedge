// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
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
	// check type
	typeExist := false
	for _, v := range nodemodule.TypesOfNode {
		if v == req.NodeType {
			typeExist = true
		}
	}

	if !typeExist {
		return nil, nodemodule.ErrNodeTypeNotExist
	}

	// get cluster
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		logger.Infof("CreateEdgeNode, can`t find cluster %d", userInfo.ClusterId)
		return nil, cluster.ErrFindCluster
	}
	uniqName := clu.GetUniqueName(cluster.ResourceNode)

	node := &nodeentity.NodeBasicInfo{
		ClusterId:        userInfo.ClusterId,
		GroupId:          userInfo.GroupId,
		UserId:           userInfo.UserId,
		NodeName:         req.Name,
		NodeType:         req.NodeType,
		UniqueName:       uniqName,
		Arch:             "",
		CpuCore:          0,
		Memory:           0,
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

//////// node type ////////
func ListEdgeNodeType(userInfo proto.ApulisHeader, req *nodemodule.ListEdgeNodeTypeReq) ([]string, error) {
	return nodemodule.TypesOfNode, nil
}

//////// arch type ////////
func ListArchType(userInfo proto.ApulisHeader, req *nodemodule.ListArchTypeReq) ([]string, error) {
	return []string{cluster.ArchArm, cluster.ArchX86}, nil
}
