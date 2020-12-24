// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageentity

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	"time"
)

// table contants
const (
	TableContainerImageOrg string = "ContainerImageOrgs"
)

// The simplest implement, add RBAC if need
type ContainerImageOrg struct {
	ID           int64     `gorm:"column:Id;primary_key"                                    json:"id"`
	ClusterId    int64     `gorm:"uniqueIndex:cluster_org;column:ClusterId;not null"        json:"clusterId"`
	OrgName      string    `gorm:"uniqueIndex:cluster_org;column:OrgName;size:255;not null" json:"orgName"`
	OwnerGroupId int64     `gorm:"column:OwnerGroupId;not null"                             json:"ownerGroupId"`
	OwnerUserId  int64     `gorm:"column:OwnerUserId;not null"                              json:"ownerUserId"`
	CreateAt     time.Time `gorm:"column:CreateAt;not null"                                 json:"createAt"`
	UpdateAt     time.Time `gorm:"column:UpdateAt;not null"                                 json:"updateAt"`
}

func (ContainerImageOrg) TableName() string {
	return TableContainerImageOrg
}

func CreateImageOrg(orgInfo *ContainerImageOrg) error {
	return apulisdb.Db.Create(orgInfo).Error
}

func UpdateImageOrg(orgInfo *ContainerImageOrg) error {
	return apulisdb.Db.Save(orgInfo).Error
}

func DeleteImageOrg(orgInfo *ContainerImageOrg) error {
	return apulisdb.Db.Delete(orgInfo).Error
}
