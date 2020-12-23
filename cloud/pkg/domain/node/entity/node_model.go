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

// NodeName is unique in Cluster/Group/User
type NodeBasicInfo struct {
	ID               int64     `gorm:"column:Id;primary_key"                                   json:"id" binding:"required"`
	NodeName         string    `gorm:"uniqueIndex:user_node;column:NodeName;size:255;not null" json:"name" binding:"required"`
	ClusterId        int64     `gorm:"uniqueIndex:user_node;column:ClusterId;not null"         json:"clusterId" binding:"required"`
	GroupId          int64     `gorm:"uniqueIndex:user_node;column:GroupId;not null"           json:"groupId" binding:"required"`
	UserId           int64     `gorm:"uniqueIndex:user_node;column:UserId;not null"            json:"userId" binding:"required"`
	Status           string    `gorm:"column:Status;size:255;not null"                         json:"status" binding:"required"`
	Roles            string    `gorm:"column:Roles;size:255;not null"                          json:"roles" binding:"required"`
	ContainerRuntime string    `gorm:"column:ContainerRuntime;size:255;not null"               json:"runtime" binding:"required"`
	OsImage          string    `gorm:"column:OsImage;size:255;not null"                        json:"osImage" binding:"required"`
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

func DeleteNode(nodeInfo *NodeBasicInfo) error {
	return apulisdb.Db.Delete(nodeInfo).Error
}
