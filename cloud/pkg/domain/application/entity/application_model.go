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
	ID               int64     `gorm:"column:Id;primary_key"                                 json:"id"`
	AppName          string    `gorm:"uniqueIndex:user_app;column:AppName;size:255;not null" json:"appName"`
	ClusterId        int64     `gorm:"uniqueIndex:user_app;column:ClusterId;not null"        json:"clusterId"`
	GroupId          int64     `gorm:"uniqueIndex:user_app;column:GroupId;not null"          json:"groupId"`
	UserId           int64     `gorm:"uniqueIndex:user_app;column:UserId;not null"           json:"userId"`
	AppType          string    `gorm:"column:appType;size:128;not null"                      json:"appType"`
	FunctionType     string    `gorm:"column:FunctionType;size:1024;not null"                json:"functionType"`
	Description      string    `gorm:"column:Description;size:1024;not null"                 json:"description"`
	LatestPubVersion string    `gorm:"column:LatestPubVersion;size:255;not null"             json:"latestPubVersion"`
	CreateAt         time.Time `gorm:"column:CreateAt;not null"                              json:"createAt"`
	UpdateAt         time.Time `gorm:"column:UpdateAt;not null"                              json:"updateAt"`
}

type ApplicationVersionInfo struct {
	ID                    int64     `gorm:"column:Id;primary_key"                                    json:"id"`
	AppName               string    `gorm:"uniqueIndex:app_version;column:AppName;size:255;not null" json:"appName"`
	ClusterId             int64     `gorm:"uniqueIndex:app_version;column:ClusterId;not null"        json:"clusterId"`
	GroupId               int64     `gorm:"uniqueIndex:app_version;column:GroupId;not null"          json:"groupId"`
	UserId                int64     `gorm:"uniqueIndex:app_version;column:UserId;not null"           json:"userId"`
	Version               string    `gorm:"uniqueIndex:app_version;column:Version;size:255;not null" json:"version"`
	Status                string    `gorm:"column:Status;size:255;not null"                          json:"status"`
	ArchType              string    `gorm:"column:ArchType;size:64;not null"                         json:"archType"`
	ContainerImage        string    `gorm:"column:containerImage;size:255;not null"                  json:"containerImage"`
	ContainerImageVersion string    `gorm:"column:containerImageVersion;size:255;not null"           json:"containerImageVersion"`
	ContainerImagePath    string    `gorm:"column:containerImagePath;size:255;not null"              json:"containerImagePath"`
	CpuQuota              float32   `gorm:"column:CpuQuota;not null"                                 json:"cpuQuota"`
	MaxCpuQuota           float32   `gorm:"column:MaxCpuQuota;not null"                              json:"maxCpuQuota"`
	MemoryQuota           float32   `gorm:"column:MemoryQuota;not null"                              json:"memoryQuota"`
	MaxMemoryQuota        float32   `gorm:"column:MaxMemoryQuota;not null"                           json:"MaxMemoryQuota"`
	RestartPolicy         string    `gorm:"column:RestartPolicy;size:128;not null"                   json:"restartPolicy"`
	Network               string    `gorm:"column:Network;size:512;not null"                         json:"network"`
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
