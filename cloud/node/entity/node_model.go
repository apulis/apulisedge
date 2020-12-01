// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/database"
)

type NodeBasicInfo struct {
	ID               int64  `gorm:"primary_key" json:"id" binding:"required"`
	Name             string `gorm:"not null" json:"name" binding:"required"`
	Status           string `gorm:"not null" json:"status" binding:"required"`
	Roles            string `gorm:"not null" json:"roles" binding:"required"`
	ContainerRuntime string `gorm:"not null" json:"runtime" binding:"required"`
	OsImage          string `gorm:"not null" json:"osImage" binding:"required"`
	ProviderId       string `json:"providerId"`
	InterIp          string `gorm:"not null" json:"interIp"`
	OuterIp          string `json:"outerIp"`
	CreateTime       string `json:"createTime"`
}

type NodeDetailInfo struct {
	Name             string `json:"name" binding:"required"`
	Status           string `json:"status" binding:"required"`
	Roles            string `json:"roles" binding:"required"`
	ContainerRuntime string `json:"runtime" binding:"required"`
	OsImage          string `json:"osImage" binding:"required"`
	ProviderId       string `json:"providerId"`
	InterIp          string `json:"interIp"`
	OuterIp          string `json:"outerIp"`
	CreateTime       string `json:"createTime"`
}

func CreateNode(nodeInfo *NodeBasicInfo) error {
	return apulisdb.Db.Create(nodeInfo).Error
}
