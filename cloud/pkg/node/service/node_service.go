// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"fmt"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/node/entity"
	"time"
)

var logger = loggers.LogInstance()

func CreateEdgeNode(userId int64, userName string, name string) (*nodeentity.NodeBasicInfo, error) {
	node := &nodeentity.NodeBasicInfo{
		UserId:           userId,
		UserName:         userName,
		Name:             name,
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

	whereQueryStr := fmt.Sprintf("UserId = '%s' and Name = '%s'", userId, nodeName)
	res := apulisdb.Db.Where(whereQueryStr).First(&nodeInfo)

	if res.Error != nil {
		return &nodeInfo, res.Error
	}

	return &nodeInfo, nil
}
