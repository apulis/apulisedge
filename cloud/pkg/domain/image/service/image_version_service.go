// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
)

// list container image version
func ListContainerImageVersion(req *imagemodule.ListContainerImageVersionReq) (*[]imageentity.UserContainerImageVersionInfo, error) {
	var imageVerInfos []imageentity.UserContainerImageVersionInfo

	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ?",
			req.ClusterId, req.GroupId, req.UserId, req.ImageName, req.OrgName).
		Find(&imageVerInfos)

	if res.Error != nil {
		return &imageVerInfos, res.Error
	}

	return &imageVerInfos, nil
}

// delete container image version
func DeleteContainterImageVersion(req *imagemodule.DeleteContainerImageVersionReq) error {
	var imageVerInfo imageentity.UserContainerImageVersionInfo

	// get image and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ? and ImageVersion = ?",
			req.ClusterId, req.GroupId, req.UserId, req.ImageName, req.OrgName, req.ImageVersion).
		First(&imageVerInfo)
	if res.Error != nil {
		return res.Error
	}

	return imageentity.DeleteContainerImageVersion(&imageVerInfo)
}
