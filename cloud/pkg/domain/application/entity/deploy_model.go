// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableApplicationDeployInfo string = "ApplicationDeployInfos"
)

type ApplicationDeployInfo struct {
	ID                    int64     `gorm:"column:Id;primary_key"                                         json:"id" binding:"required"`
	AppName               string    `gorm:"uniqueIndex:user_app;index:;column:AppName;size:255;not null"  json:"appName" binding:"required"`
	UserId                int64     `gorm:"uniqueIndex:user_app;column:UserId;not null"                   json:"userId" binding:"required"`
	UserName              string    `gorm:"column:UserName;size:255;not null"                             json:"userName" binding:"required"`
	NodeName              string    `gorm:"uniqueIndex:user_app;column:NodeName;size:255;not null"        json:"nodeName" binding:"required"`
	Status                string    `gorm:"column:Status;size:255;not null"                               json:"status" binding:"required"`
	ArchType              int       `gorm:"column:ArchType;not null"                                      json:"archType" binding:"required"`
	Version               string    `gorm:"column:Version;size:255;not null"                              json:"version" binding:"required"`
	Namespace             string    `gorm:"column:Namespace;size:255;not null"                            json:"namespace" binding:"required"`
	ContainerImage        string    `gorm:"column:containerImage;size:255;not null"                       json:"containerImage" binding:"required"`
	ContainerImageVersion string    `gorm:"column:containerImageVersion;size:255;not null"                json:"containerImageVersion" binding:"required"`
	ContainerImagePath    string    `gorm:"column:containerImagePath;size:255;not null"                   json:"containerImagePath" binding:"required"`
	CpuQuota              float32   `gorm:"column:CpuQuota;not null"                                      json:"CpuQuota" binding:"required"`
	MemoryQuota           int       `gorm:"column:MemoryQuota;not null"                                   json:"MemoryQuota" binding:"required"`
	ContainerPort         int       `gorm:"column:ContainerPort;not null"                                 json:"ContainerPort"`
	HostPort              int       `gorm:"column:HostPort"                                               json:"HostPort"`
	CreateAt              time.Time `gorm:"column:CreateAt;not null"                                      json:"createAt"`
	UpdateAt              time.Time `gorm:"column:UpdateAt;not null"                                      json:"updateAt"`
}

func (ApplicationDeployInfo) TableName() string {
	return TableApplicationDeployInfo
}

func CreateAppDeploy(deployInfo *ApplicationDeployInfo) error {
	return apulisdb.Db.Create(deployInfo).Error
}

func UpdateAppDeploy(deployInfo *ApplicationDeployInfo) error {
	return apulisdb.Db.Save(deployInfo).Error
}

// TODO delete deploy
