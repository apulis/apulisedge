// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// list container org
func ListContainerImageOrg(userInfo proto.ApulisHeader, req *imagemodule.ListContainerImageOrgReq) (*[]imageentity.ContainerImageOrg, error) {
	var imageOrgs []imageentity.ContainerImageOrg

	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and OrgName = ? and OwnerGroupId = ? and OwnerUserId = ?",
			userInfo.ClusterId, req.OrgName, userInfo.GroupId, userInfo.UserId).
		Find(&imageOrgs)

	if res.Error != nil {
		return &imageOrgs, res.Error
	}

	return &imageOrgs, nil
}

// delete container org
func DeleteContainterImageOrg(userInfo proto.ApulisHeader, req *imagemodule.DeleteContainerImageOrgReq) error {

	// check if org have images
	var total int64
	apulisdb.Db.Model(&imageentity.UserContainerImageInfo{}).
		Where("ClusterId = ? and OrgName = ?", userInfo.ClusterId, req.OrgName).
		Count(&total)
	if total != 0 {
		return imagemodule.ErrOrgImageNotEmpty
	}

	// delete org
	res := apulisdb.Db.Where("ClusterId = ? and OrgName = ?", userInfo.ClusterId, req.OrgName).
		Delete(imageentity.ContainerImageOrg{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
