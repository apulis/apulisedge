// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableNodeGroupInfo string = "NodeGroupInfos"
)

// NodeGroup is unique in Cluster/Group/User
type NodeGroupInfo struct {
	ID        int64     `gorm:"column:Id;primary_key"                                     json:"id"`
	ClusterId int64     `gorm:"uniqueIndex:node_group;column:ClusterId;not null"          json:"clusterId"`
	GroupId   int64     `gorm:"uniqueIndex:node_group;column:GroupId;not null"            json:"groupId"`
	UserId    int64     `gorm:"uniqueIndex:node_group;column:UserId;not null"             json:"userId"`
	GroupName string    `gorm:"uniqueIndex:node_group;column:GroupName;size:255;not null" json:"groupName"`
	CreateAt  time.Time `gorm:"column:CreateAt;not null"                                  json:"createAt"`
	UpdateAt  time.Time `gorm:"column:UpdateAt;not null"                                  json:"updateAt"`
}

func (NodeGroupInfo) TableName() string {
	return TableNodeGroupInfo
}

func CreateGroup(groupInfo *NodeGroupInfo) error {
	return apulisdb.Db.Create(groupInfo).Error
}

func UpdateGroup(groupInfo *NodeGroupInfo) error {
	return apulisdb.Db.Save(groupInfo).Error
}

func DeleteGroup(groupInfo *NodeGroupInfo) error {
	return apulisdb.Db.Delete(groupInfo).Error
}
