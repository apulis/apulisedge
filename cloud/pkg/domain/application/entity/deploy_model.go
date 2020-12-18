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
	ID         int64     `gorm:"column:Id;primary_key"                                          json:"id" binding:"required"`
	ClusterId  int64     `gorm:"uniqueIndex:app_deploy;column:ClusterId;not null"               json:"clusterId" binding:"required"`
	GroupId    int64     `gorm:"uniqueIndex:app_deploy;column:GroupId;not null"                 json:"groupId" binding:"required"`
	UserId     int64     `gorm:"uniqueIndex:app_deploy;column:UserId;not null"                  json:"userId" binding:"required"`
	NodeName   string    `gorm:"uniqueIndex:app_deploy;column:NodeName;size:255;not null"       json:"nodeName" binding:"required"`
	AppName    string    `gorm:"uniqueIndex:app_deploy;index:;column:AppName;size:255;not null" json:"appName" binding:"required"`
	Version    string    `gorm:"column:Version;size:255;not null"                               json:"version" binding:"required"`
	Status     string    `gorm:"column:Status;size:255;not null"                                json:"status" binding:"required"`
	DeployUUID string    `gorm:"column:DeployUUID;index:uuid;size:255;not null"                 json:"deployUUID" binding:"required"`
	CreateAt   time.Time `gorm:"column:CreateAt;not null"                                       json:"createAt"`
	UpdateAt   time.Time `gorm:"column:UpdateAt;not null"                                       json:"updateAt"`
}

type PortMappingInfo struct {
	ContainerPort int `json:"containerPort"`
	HostPort      int `json:"hostPort"`
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

func DeleteAppDeploy(deployInfo *ApplicationDeployInfo) error {
	return apulisdb.Db.Delete(deployInfo).Error
}
