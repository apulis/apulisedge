// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeentity

import (
	"time"

	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
)

// table contants
const (
	TableNodeBasicInfo string = "NodeBasicInfos"
)

// NodeName is unique in Cluster/Group/User
type NodeBasicInfo struct {
	ID               int64     `gorm:"column:Id;primary_key"                                   json:"id"`
	NodeName         string    `gorm:"uniqueIndex:user_node;column:NodeName;size:255;not null" json:"name"`
	ClusterId        int64     `gorm:"uniqueIndex:user_node;column:ClusterId;not null"         json:"clusterId"`
	GroupId          int64     `gorm:"uniqueIndex:user_node;column:GroupId;not null"           json:"groupId"`
	UserId           int64     `gorm:"uniqueIndex:user_node;column:UserId;not null"            json:"userId"`
	NodeType         string    `gorm:"column:NodeType;size:128;not null"                       json:"nodeType"`
	Arch             string    `gorm:"column:Arch;size:128;not null"                           json:"arch"`
	UniqueName       string    `gorm:"column:UniqueName;size:255;not null"                     json:"uniqueName"`
	CpuCore          int       `gorm:"column:CpuCore;not null"                                 json:"cpuCore"`
	Memory           int64     `gorm:"column:Memory;not null"                                  json:"memory"`
	Status           string    `gorm:"column:Status;size:255;not null"                         json:"status"`
	Roles            string    `gorm:"column:Roles;size:255;not null"                          json:"roles"`
	ContainerRuntime string    `gorm:"column:ContainerRuntime;size:255;not null"               json:"runtime"`
	OsImage          string    `gorm:"column:OsImage;size:255;not null"                        json:"osImage"`
	InterIp          string    `gorm:"column:InterIp;size:255;not null"                        json:"interIp"`
	OuterIp          string    `gorm:"column:OuterIp;size:255"                                 json:"outerIp"`
	CreateAt         time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt         time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
}

// NodeOfBatchInfo is one case of batch before comfirmation
type NodeOfBatchInfo struct {
	ID       int64     `gorm:"column:Id;primary_key"                                   json:"id"`
	NodeName string    `gorm:"uniqueIndex:user_node;column:NodeName;size:255;not null" json:"name"`
	Address  string    `gorm:"column:Address" json:"address"`
	Port     string    `gorm:"column:Port" json:"port"`
	Password string    `gorm:"column:Password" json:"password"`
	CreateAt time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
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
