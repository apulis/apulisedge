// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableUserContainerImageInfo        string = "UserContainerImageInfos"
	TableUserContainerImageVersionInfo string = "UserContainerImageVersionInfos"
)

type UserContainerImageInfo struct {
	ID        int64     `gorm:"column:Id;primary_key"                                   json:"id"`
	ClusterId int64     `gorm:"uniqueIndex:user_img;column:ClusterId;not null"          json:"clusterId"`
	GroupId   int64     `gorm:"uniqueIndex:user_img;column:GroupId;not null"            json:"groupId"`
	UserId    int64     `gorm:"uniqueIndex:user_img;column:UserId;not null"             json:"userId"`
	ImageName string    `gorm:"uniqueIndex:user_img;column:ImageName;size:255;not null" json:"imageName"`
	OrgName   string    `gorm:"uniqueIndex:user_img;column:OrgName;size:255;not null"   json:"orgName"`
	CreateAt  time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt  time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
}

type UserContainerImageVersionInfo struct {
	ID              int64     `gorm:"column:Id;primary_key"                                              json:"id"`
	ClusterId       int64     `gorm:"uniqueIndex:user_img_version;column:ClusterId;not null"             json:"clusterId"`
	GroupId         int64     `gorm:"uniqueIndex:user_img_version;column:GroupId;not null"               json:"groupId"`
	UserId          int64     `gorm:"uniqueIndex:user_img_version;column:UserId;not null"                json:"userId"`
	ImageName       string    `gorm:"uniqueIndex:user_img_version;column:ImageName;size:255;not null"    json:"imageName"`
	OrgName         string    `gorm:"uniqueIndex:user_img_version;column:OrgName;size:255;not null"      json:"orgName"`
	ImageId         string    `gorm:"column:ImageId;size:512;not null"                                   json:"imageId"`
	ImageVersion    string    `gorm:"uniqueIndex:user_img_version;column:ImageVersion;size:255;not null" json:"imageVersion"`
	ImageSize       float32   `gorm:"column:ImageSize;not null"                                          json:"imageSize"`
	DownloadCommand string    `gorm:"column:DownloadCommand;size:512;not null"                           json:"DownloadCommand"`
	CreateAt        time.Time `gorm:"column:CreateAt;not null"                                           json:"createAt"`
	UpdateAt        time.Time `gorm:"column:UpdateAt;not null"                                           json:"updateAt"`
}

// image list
func (UserContainerImageInfo) TableName() string {
	return TableUserContainerImageInfo
}

func CreateContainerImage(imgInfo *UserContainerImageInfo) error {
	return apulisdb.Db.Create(imgInfo).Error
}

func UpdateContainerImage(imgInfo *UserContainerImageInfo) error {
	return apulisdb.Db.Save(imgInfo).Error
}

func DeleteContainerImage(imgInfo *UserContainerImageInfo) error {
	return apulisdb.Db.Delete(imgInfo).Error
}

// image version
func (UserContainerImageVersionInfo) TableName() string {
	return TableUserContainerImageVersionInfo
}

func CreateContainerImageVersion(imgVerInfo *UserContainerImageVersionInfo) error {
	return apulisdb.Db.Create(imgVerInfo).Error
}

func UpdateContainerImageVersion(imgVerInfo *UserContainerImageVersionInfo) error {
	return apulisdb.Db.Save(imgVerInfo).Error
}

func DeleteContainerImageVersion(imgVerInfo *UserContainerImageVersionInfo) error {
	return apulisdb.Db.Delete(imgVerInfo).Error
}
