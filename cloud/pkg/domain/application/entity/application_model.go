// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableApplicationBasicInfo   string = "ApplicationBasicInfos"
	TableApplicationVersionInfo string = "ApplicationVersionInfos"
)

type ApplicationBasicInfo struct {
	ID               int64     `gorm:"column:Id;primary_key"                                 json:"id" binding:"required"`
	AppName          string    `gorm:"uniqueIndex:user_app;column:AppName;size:255;not null" json:"appName" binding:"required"`
	ClusterId        int64     `gorm:"uniqueIndex:user_app;column:ClusterId;not null"        json:"clusterId" binding:"required"`
	GroupId          int64     `gorm:"uniqueIndex:user_app;column:GroupId;not null"          json:"groupId" binding:"required"`
	UserId           int64     `gorm:"uniqueIndex:user_app;column:UserId;not null"           json:"userId" binding:"required"`
	AppType          string    `gorm:"column:appType;size:128;not null"                      json:"appType" binding:"required"`
	FunctionType     string    `gorm:"column:FunctionType;size:1024;not null"                json:"functionType" binding:"required"`
	Description      string    `gorm:"column:Description;size:1024;not null"                 json:"description"`
	LatestPubVersion string    `gorm:"column:LatestPubVersion;size:255;not null"             json:"latestPubVersion" binding:"required"`
	CreateAt         time.Time `gorm:"column:CreateAt;not null"                              json:"createAt"`
	UpdateAt         time.Time `gorm:"column:UpdateAt;not null"                              json:"updateAt"`
}

// TODO port mapping
type ApplicationVersionInfo struct {
	ID                    int64     `gorm:"column:Id;primary_key"                                    json:"id" binding:"required"`
	AppName               string    `gorm:"uniqueIndex:app_version;column:AppName;size:255;not null" json:"appName" binding:"required"`
	ClusterId             int64     `gorm:"uniqueIndex:app_version;column:ClusterId;not null"        json:"clusterId" binding:"required"`
	GroupId               int64     `gorm:"uniqueIndex:app_version;column:GroupId;not null"          json:"groupId" binding:"required"`
	UserId                int64     `gorm:"uniqueIndex:app_version;column:UserId;not null"           json:"userId" binding:"required"`
	Version               string    `gorm:"uniqueIndex:app_version;column:Version;size:255;not null" json:"version" binding:"required"`
	Status                string    `gorm:"column:Status;size:255;not null"                          json:"status" binding:"required"`
	ArchType              string    `gorm:"column:ArchType;size:64;not null"                         json:"archType" binding:"required"`
	ContainerImage        string    `gorm:"column:containerImage;size:255;not null"                  json:"containerImage" binding:"required"`
	ContainerImageVersion string    `gorm:"column:containerImageVersion;size:255;not null"           json:"containerImageVersion" binding:"required"`
	ContainerImagePath    string    `gorm:"column:containerImagePath;size:255;not null"              json:"containerImagePath" binding:"required"`
	CpuQuota              float32   `gorm:"column:CpuQuota;not null"                                 json:"cpuQuota" binding:"required"`
	MaxCpuQuota           float32   `gorm:"column:MaxCpuQuota;not null"                              json:"maxCpuQuota" binding:"required"`
	MemoryQuota           float32   `gorm:"column:MemoryQuota;not null"                              json:"memoryQuota" binding:"required"`
	MaxMemoryQuota        float32   `gorm:"column:MaxMemoryQuota;not null"                           json:"MaxMemoryQuota" binding:"required"`
	RestartPolicy         string    `gorm:"column:RestartPolicy;size:128;not null"                   json:"restartPolicy" binding:"required"`
	Network               string    `gorm:"column:Network;size:512;not null"                         json:"network" binding:"required"`
	CreateAt              time.Time `gorm:"column:CreateAt;not null"                                 json:"createAt"`
	UpdateAt              time.Time `gorm:"column:UpdateAt;not null"                                 json:"updateAt"`
	PublishAt             string    `gorm:"column:PublishAt"                                         json:"publishAt"`
	OfflineAt             string    `gorm:"column:OfflineAt"                                         json:"offlineAt"`
}

func (ApplicationBasicInfo) TableName() string {
	return TableApplicationBasicInfo
}

func (ApplicationVersionInfo) TableName() string {
	return TableApplicationVersionInfo
}

func CreateApplication(appInfo *ApplicationBasicInfo) error {
	return apulisdb.Db.Create(appInfo).Error
}

func UpdateApplication(appInfo *ApplicationBasicInfo) error {
	return apulisdb.Db.Save(appInfo).Error
}

func DeleteApplication(appInfo *ApplicationBasicInfo) error {
	return apulisdb.Db.Delete(appInfo).Error
}

func CreateApplicationVersion(appInfo *ApplicationVersionInfo) error {
	return apulisdb.Db.Create(appInfo).Error
}

func UpdateApplicationVersion(appInfo *ApplicationVersionInfo) error {
	return apulisdb.Db.Save(appInfo).Error
}

func DeleteApplicationVersion(appInfo *ApplicationVersionInfo) error {
	return apulisdb.Db.Delete(appInfo).Error
}
