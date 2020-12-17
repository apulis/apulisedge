// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"time"
)

// create container org
func CreateContainerImageOrg(userInfo proto.ApulisHeader, req *imagemodule.CreateContainerImageOrgReq) (*imageentity.ContainerImageOrg, error) {
	orgInfo := &imageentity.ContainerImageOrg{
		ClusterId:    userInfo.ClusterId,
		OrgName:      req.OrgName,
		OwnerGroupId: userInfo.GroupId,
		OwnerUserId:  userInfo.UserId,
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}

	err := imageentity.CreateImageOrg(orgInfo)
	if err != nil {
		return nil, err
	}

	return orgInfo, nil
}

// list container org
func ListContainerImageOrg(userInfo proto.ApulisHeader, req *imagemodule.ListContainerImageOrgReq) (*[]imageentity.ContainerImageOrg, int, error) {
	var imageOrgs []imageentity.ContainerImageOrg

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and OwnerGroupId = ? and OwnerUserId = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Find(&imageOrgs)

	if res.Error != nil {
		return &imageOrgs, total, res.Error
	}

	total = int(res.RowsAffected)
	return &imageOrgs, total, nil
}

// delete container org
func DeleteContainterImageOrg(userInfo proto.ApulisHeader, req *imagemodule.DeleteContainerImageOrgReq) error {
	var orgInfo imageentity.ContainerImageOrg

	// check if org have images
	var total int64
	apulisdb.Db.Model(&imageentity.UserContainerImageInfo{}).
		Where("ClusterId = ? and OrgName = ?", userInfo.ClusterId, req.OrgName).
		Count(&total)
	if total != 0 {
		return imagemodule.ErrOrgImageNotEmpty
	}

	// delete org
	res := apulisdb.Db.
		Where("ClusterId = ? and OrgName = ? and OwnerGroupId = ? and OwnerUserId = ?",
			userInfo.ClusterId, req.OrgName, userInfo.GroupId, userInfo.UserId).
		First(&orgInfo)
	if res.Error != nil {
		return res.Error
	}

	res = apulisdb.Db.Where("ClusterId = ? and OrgName = ? and OwnerGroupId = ? and OwnerUserId = ?",
		userInfo.ClusterId, req.OrgName, userInfo.GroupId, userInfo.UserId).
		Delete(imageentity.ContainerImageOrg{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
