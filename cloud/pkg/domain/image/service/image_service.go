// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

var logger = loggers.LogInstance()

// list container images
func ListContainerImage(userInfo proto.ApulisHeader, req *imagemodule.ListContainerImageReq) ([]imageentity.UserContainerImageInfo, int, error) {
	var imageInfos []imageentity.UserContainerImageInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Group("ImageName").
		Group("OrgName").
		Select("ImageName, OrgName", "UpdateAt").
		Find(&imageInfos)

	if res.Error != nil {
		return imageInfos, total, res.Error
	}

	return imageInfos, int(res.RowsAffected), nil
}

// delete container images
func DeleteContainerImage(userInfo proto.ApulisHeader, req *imagemodule.DeleteContainerImageReq) error {
	var imageInfo imageentity.UserContainerImageInfo

	// check if any image version exist
	var total int64
	apulisdb.Db.Model(&imageentity.UserContainerImageVersionInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.ImageName, req.OrgName).
		Count(&total)
	if total != 0 {
		return imagemodule.ErrImageVersionExist
	}

	// get image and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ImageName = ? and OrgName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.ImageName, req.OrgName).
		First(&imageInfo)
	if res.Error != nil {
		return res.Error
	}

	return imageentity.DeleteContainerImage(&imageInfo)
}
