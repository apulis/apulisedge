// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableUserContainerImageInfo string = "UserContainerImageInfos"
)

type UserContainerImageInfo struct {
	ID        int64     `gorm:"column:Id;primary_key"                                   json:"id" binding:"required"`
	ClusterId int64     `gorm:"uniqueIndex:user_img;column:ClusterId;not null"          json:"clusterId" binding:"required"`
	GroupId   int64     `gorm:"uniqueIndex:user_img;column:GroupId;not null"            json:"groupId" binding:"required"`
	UserId    int64     `gorm:"uniqueIndex:user_img;column:UserId;not null"             json:"userId" binding:"required"`
	ImageName string    `gorm:"uniqueIndex:user_img;column:ImageName;size:255;not null" json:"imageName" binding:"required"`
	OrgName   string    `gorm:"uniqueIndex:user_img;column:OrgName;size:255;not null"   json:"orgName" binding:"required"`
	CreateAt  time.Time `gorm:"column:CreateAt;not null"                                json:"createAt"`
	UpdateAt  time.Time `gorm:"column:UpdateAt;not null"                                json:"updateAt"`
}

//type UserContainerImageInfo struct {
//	ID           int64     `gorm:"column:Id;primary_key"                                       json:"id" binding:"required"`
//	ClusterId    int64     `gorm:"uniqueIndex:user_img;column:ClusterId;not null"              json:"clusterId" binding:"required"`
//	GroupId      int64     `gorm:"uniqueIndex:user_img;column:GroupId;not null"                json:"groupId" binding:"required"`
//	UserId       int64     `gorm:"uniqueIndex:user_img;column:UserId;not null"                 json:"userId" binding:"required"`
//	ImageName    string    `gorm:"uniqueIndex:user_img;column:ImageName;size:255;not null"     json:"imageName" binding:"required"`
//	ImageId      string    `gorm:"column:ImageId;size:512;not null"                            json:"imageId" binding:"required"`
//	OrgName      string    `gorm:"uniqueIndex:user_img;column:OrgName;size:255;not null"       json:"orgName" binding:"required"`
//	ImageVersion string    `gorm:"uniqueIndex:user_img;column:ImageVersion;size:255;not null"  json:"imageVersion" binding:"required"`
//	ImageSize    float32   `gorm:"column:ImageSize;size:255;not null"                          json:"imageSize" binding:"required"`
//	CreateAt     time.Time `gorm:"column:CreateAt;not null"                                    json:"createAt"`
//	UpdateAt     time.Time `gorm:"column:UpdateAt;not null"                                    json:"updateAt"`
//}

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
