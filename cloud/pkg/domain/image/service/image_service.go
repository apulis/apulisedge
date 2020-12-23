// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"gorm.io/gorm"
	"time"
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
		Find(&imageInfos)

	if res.Error != nil {
		return imageInfos, total, res.Error
	}

	return imageInfos, int(res.RowsAffected), nil
}

// add container images
func AddContainerImage(userInfo proto.ApulisHeader, orgName string, imgName string, imgVer string, imgRepo string) error {
	var tmpImgInfo imageentity.UserContainerImageInfo
	var tmpImgVerInfo imageentity.UserContainerImageVersionInfo

	var imgExsit bool
	var imgVerExsit bool

	// check image
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and OrgName = ? and ImageName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, orgName, imgName).
		First(&tmpImgInfo)
	if res.Error == nil {
		logger.Errorf("AddContainerImage image exist, org name = %s, img name = %s", orgName, imgName)
		imgExsit = true
	} else if res.Error == gorm.ErrRecordNotFound {
		imgExsit = false
	} else {
		logger.Errorf("AddContainerImage get img failed. err = %v", res.Error)
		return res.Error
	}

	// check image version
	res = apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and OrgName = ? and ImageName = ? and ImageVersion = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, orgName, imgName, imgVer).
		First(&tmpImgVerInfo)
	if res.Error == nil {
		logger.Errorf("AddContainerImage image version exist, org name = %s, img name = %s, ver = %s", orgName, imgName, imgVer)
		imgVerExsit = true
	} else if res.Error == gorm.ErrRecordNotFound {
		imgVerExsit = false
	} else {
		logger.Errorf("AddContainerImage get img version failed. err = %v", res.Error)
		return res.Error
	}

	if !imgExsit {
		imageInfo := imageentity.UserContainerImageInfo{
			ClusterId: userInfo.ClusterId,
			GroupId:   userInfo.GroupId,
			UserId:    userInfo.UserId,
			ImageName: imgName,
			OrgName:   orgName,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		}

		err := imageentity.CreateContainerImage(&imageInfo)
		if err != nil {
			logger.Errorf("AddContainerImage create image failed. err = %v", err)
			return err
		}
	} else if imgExsit { // just update timestamp
		tmpImgInfo.UpdateAt = time.Now()
		_ = imageentity.UpdateContainerImage(&tmpImgInfo)
	}

	if !imgVerExsit {
		imageVerInfo := imageentity.UserContainerImageVersionInfo{
			ClusterId:       userInfo.ClusterId,
			GroupId:         userInfo.GroupId,
			UserId:          userInfo.UserId,
			ImageName:       imgName,
			OrgName:         orgName,
			ImageId:         "",
			ImageVersion:    imgVer,
			ImageSize:       float32(0),
			DownloadCommand: imagemodule.DockerPullPrefix + imagemodule.BlankString + imgRepo,
			CreateAt:        time.Now(),
			UpdateAt:        time.Now(),
		}

		err := imageentity.CreateContainerImageVersion(&imageVerInfo)
		if err != nil {
			logger.Errorf("AddContainerImage create image version failed. err = %v", err)
			return err
		}
	} else if imgVerExsit { // just update timestamp
		tmpImgVerInfo.UpdateAt = time.Now()
		_ = imageentity.UpdateContainerImageVersion(&tmpImgVerInfo)
	}

	return nil
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

	err := imageentity.DeleteContainerImage(&imageInfo)
	if err != nil {
		return err
	}

	return nil
}
