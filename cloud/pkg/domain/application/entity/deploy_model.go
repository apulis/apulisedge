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
	ID         int64     `gorm:"column:Id;primary_key"                                          json:"id"`
	ClusterId  int64     `gorm:"uniqueIndex:app_deploy;column:ClusterId;not null"               json:"clusterId"`
	GroupId    int64     `gorm:"uniqueIndex:app_deploy;column:GroupId;not null"                 json:"groupId"`
	UserId     int64     `gorm:"uniqueIndex:app_deploy;column:UserId;not null"                  json:"userId"`
	NodeName   string    `gorm:"uniqueIndex:app_deploy;column:NodeName;size:255;not null"       json:"nodeName"`
	AppName    string    `gorm:"uniqueIndex:app_deploy;index:;column:AppName;size:255;not null" json:"appName"`
	Version    string    `gorm:"column:Version;size:255;not null"                               json:"version"`
	Status     string    `gorm:"column:Status;size:255;not null"                                json:"status"`
	DeployUUID string    `gorm:"column:DeployUUID;index:uuid;size:255;not null"                 json:"deployUUID"`
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
