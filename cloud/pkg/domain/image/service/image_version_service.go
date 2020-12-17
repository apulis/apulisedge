// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// list container image version
func ListContainerImageVersion(userInfo proto.ApulisHeader, req *imagemodule.ListContainerImageVersionReq) (*[]imageentity.UserContainerImageVersionInfo, int, error) {
	var imageVerInfos []imageentity.UserContainerImageVersionInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.ImageName, req.OrgName).
		Find(&imageVerInfos)

	if res.Error != nil {
		return &imageVerInfos, total, res.Error
	}

	total = int(res.RowsAffected)
	return &imageVerInfos, total, nil
}

// delete container image version
func DeleteContainterImageVersion(userInfo proto.ApulisHeader, req *imagemodule.DeleteContainerImageVersionReq) error {
	var imageVerInfo imageentity.UserContainerImageVersionInfo

	// get image and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ? and ImageVersion = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.ImageName, req.OrgName, req.ImageVersion).
		First(&imageVerInfo)
	if res.Error != nil {
		return res.Error
	}

	return imageentity.DeleteContainerImageVersion(&imageVerInfo)
}
