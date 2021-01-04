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
	ID        int64     `gorm:"column:Id;primary_key"                                   json:"id"`
	NodeID    int64     `gorm:"uniqueIndex:user_node;column:NodeID;not null"         json:"NodeID"`
	ClusterId int64     `gorm:"uniqueIndex:user_node;column:ClusterId;not null"         json:"clusterId"`
	GroupId   int64     `gorm:"uniqueIndex:user_node;column:GroupId;not null"           json:"groupId"`
	UserId    int64     `gorm:"uniqueIndex:user_node;column:UserId;not null"            json:"userId"`
	NodeName  string    `gorm:"uniqueIndex:user_node;column:NodeName;size:255;not null" json:"name"`
	NodeType  string    `gorm:"column:NodeType;size:128;not null"                       json:"nodeType"`
	Arch      string    `gorm:"column:Arch;size:128;not null"                           json:"arch"`
	Address   string    `gorm:"uniqueIndex:user_node;column:Address" json:"address"`
	Port      string    `gorm:"uniqueIndex:user_node;column:Port" json:"port"`
	Sudoer    string    `gorm:"uniqueIndex:user_node;column:Sudoer"                                         json:"sudoer"`
	Password  string    `gorm:"column:Password" json:"password"`
	IsConfirm bool      `gorm:"column:IsConfirm" json:"isConfirm"`
	CreateAt  time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt  time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
}

type BatchTaskRecord struct {
	ID        int64     `gorm:"column:Id;primary_key"                                   json:"id"`
	ClusterId int64     `gorm:"column:ClusterId;not null"         json:"clusterId"`
	GroupId   int64     `gorm:"column:GroupId;not null"           json:"groupId"`
	UserId    int64     `gorm:"column:UserId;not null"            json:"userId"`
	Status    string    `gorm:"column:Status" json:"status"`
	ErrMsg    string    `gorm:"column:ErrMsg" json:"errMsg"`
	FilePath  string    `gorm:"column:FilePath" json:"filePath"`
	CreateAt  time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt  time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
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

func CreateNodeOfBatch(nodeInfo *NodeOfBatchInfo) error {
	return apulisdb.Db.Create(nodeInfo).Error
}

func DeleteNodeOfBatch(nodeInfo *NodeOfBatchInfo) error {
	return apulisdb.Db.Delete(nodeInfo).Error
}

func ConfirmNodesBatch(nodeInfo *NodeOfBatchInfo) error {
	err := apulisdb.Db.Model(nodeInfo).Update("NodeID", nodeInfo.NodeID).Error
	if err != nil {
		return err
	}
	err = apulisdb.Db.Model(nodeInfo).Update("IsConfirm", true).Error
	return err
}

func CreateBatchTask(taskInfo *BatchTaskRecord) error {
	return apulisdb.Db.Create(taskInfo).Error
}

func FinishBatchTask(taskInfo *BatchTaskRecord) error {
	return apulisdb.Db.Model(taskInfo).Update("Status", "finish").Error
}
