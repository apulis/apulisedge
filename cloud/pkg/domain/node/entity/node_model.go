// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableNodeBasicInfo string = "NodeBasicInfos"
)

type NodeBasicInfo struct {
	ID               int64     `gorm:"column:Id;primary_key"                                   json:"id" binding:"required"`
	NodeName         string    `gorm:"uniqueIndex:user_node;column:NodeName;size:255;not null" json:"name" binding:"required"`
	UserId           int64     `gorm:"uniqueIndex:user_node;column:UserId;not null"            json:"userId" binding:"required"`
	UserName         string    `gorm:"column:UserName;size:255;not null"                       json:"userName" binding:"required"`
	Status           string    `gorm:"column:Status;size:255;not null"                         json:"status" binding:"required"`
	Roles            string    `gorm:"column:Roles;size:255;not null"                          json:"roles" binding:"required"`
	ContainerRuntime string    `gorm:"column:ContainerRuntime;size:255;not null"               json:"runtime" binding:"required"`
	OsImage          string    `gorm:"column:OsImage;size:255;not null"                        json:"osImage" binding:"required"`
	ProviderId       string    `gorm:"column:ProviderId;size:255"                              json:"providerId"`
	InterIp          string    `gorm:"column:InterIp;size:255;not null"                        json:"interIp"`
	OuterIp          string    `gorm:"column:OuterIp;size:255"                                 json:"outerIp"`
	CreateAt         time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt         time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
}

func (NodeBasicInfo) TableName() string {
	return TableNodeBasicInfo
}

func CreateNode(nodeInfo *NodeBasicInfo) error {
	return apulisdb.Db.Create(nodeInfo).Error
}

func UpdateNode(nodeInfo *NodeBasicInfo) error {
	return apulisdb.Db.Save(nodeInfo).Error
}

func GetNode(userId int64, nodeName string) (*NodeBasicInfo, error) {
	nodeInfo := NodeBasicInfo{UserId: userId, NodeName: nodeName}
	res := apulisdb.Db.First(&nodeInfo)
	if res.Error != nil {
		return nil, res.Error
	}
	return &nodeInfo, nil
}

func DeleteNode(nodeInfo *NodeBasicInfo) error {
	return apulisdb.Db.Delete(nodeInfo).Error
}
