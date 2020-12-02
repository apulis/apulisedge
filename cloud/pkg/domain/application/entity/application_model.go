// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableApplicationBasicInfo string = "ApplicationBasicInfos"
)

type ApplicationBasicInfo struct {
	ID                    int64     `gorm:"column:Id;primary_key"                              json:"id" binding:"required"`
	AppName               string    `gorm:"uniqueIndex:user_app;column:AppName;size:255;not null" json:"name" binding:"required"`
	UserId                int64     `gorm:"uniqueIndex:user_app;column:UserId;not null"        json:"userId" binding:"required"`
	UserName              string    `gorm:"column:UserName;size:255;not null"                  json:"userName" binding:"required"`
	ArchType              int       `gorm:"column:ArchType;not null"                           json:"archType" binding:"required"`
	Version               string    `gorm:"column:Version;size:255;not null"                   json:"version" binding:"required"`
	ContainerImage        string    `gorm:"column:containerImage;size:255;not null"            json:"containerImage" binding:"required"`
	ContainerImageVersion string    `gorm:"column:containerImageVersion;size:255;not null"     json:"containerImageVersion" binding:"required"`
	ContainerImagePath    string    `gorm:"column:containerImagePath;size:255;not null"        json:"containerImagePath" binding:"required"`
	CpuQuota              float32   `gorm:"column:CpuQuota;not null"                           json:"CpuQuota" binding:"required"`
	MemoryQuota           int       `gorm:"column:MemoryQuota;not null"                        json:"MemoryQuota" binding:"required"`
	CreateAt              time.Time `gorm:"column:CreateAt;not null"                           json:"createAt"`
	UpdateAt              time.Time `gorm:"column:UpdateAt;not null"                           json:"updateAt"`
	PublishAt             time.Time `gorm:"column:PublishAt;not null"                          json:"publishAt"`
}

func (ApplicationBasicInfo) TableName() string {
	return TableApplicationBasicInfo
}

func CreateApplication(appInfo *ApplicationBasicInfo) error {
	return apulisdb.Db.Create(appInfo).Error
}

func UpdateApplication(appInfo *ApplicationBasicInfo) error {
	return apulisdb.Db.Save(appInfo).Error
}

// TODO delete app
