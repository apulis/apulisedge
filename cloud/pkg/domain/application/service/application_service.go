// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"encoding/json"
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"gorm.io/gorm"
	"strings"
	"time"
)

var logger = loggers.LogInstance()

// create edge application
// this interface can both create basic app and app version
func CreateEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.CreateEdgeApplicationReq) (string, string, error) {
	var tmpAppInfo appentity.ApplicationBasicInfo
	var tmpVerInfo appentity.ApplicationVersionInfo

	var appBasicInfo appentity.ApplicationBasicInfo
	var appVersionInfo appentity.ApplicationVersionInfo

	var appExist bool
	var verExist bool

	// check type
	archExist := true
	for _, v := range req.ArchType {
		if !cluster.IsArchValid(v) {
			archExist = false
		}
	}

	if !archExist {
		return "", "", cluster.ErrArchTypeNotExist
	}

	// check if i have this image
	imgPath, imgExsit := imageservice.DoIHaveTheImageVersion(userInfo, req.OrgName, req.ContainerImage, req.ContainerImageVersion)
	if !imgExsit {
		return "", "", appmodule.ErrImageVersionNotExist
	}

	// check network type
	if req.Network.Type == appmodule.NetworkTypePortMapping && len(req.Network.PortMappings) == 0 {
		return "", "", appmodule.ErrNetworkPortmappingEmpty
	}

	if len(req.Network.PortMappings) == 0 {
		req.Network.PortMappings = []appmodule.PortMapping{}
	}

	// check application exist
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		First(&tmpAppInfo)
	if res.Error == nil {
		logger.Errorf("CreateEdgeApplication application already exist, app name = %s", req.AppName)
		appExist = true
	} else if res.Error == gorm.ErrRecordNotFound {
		appExist = false
	} else {
		logger.Errorf("CreateEdgeApplication get application failed. err = %v", res.Error)
		return "", "", res.Error
	}

	// check version exist
	res = apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&tmpVerInfo)
	if res.Error == nil {
		logger.Errorf("CreateEdgeApplication version already exist, app name = %s, version = %s", req.AppName, req.Version)
		verExist = true
		return "", "", appmodule.ErrApplicationVersionExist
	} else if res.Error == gorm.ErrRecordNotFound {
		verExist = false
	} else {
		logger.Errorf("CreateEdgeApplication get version failed. err = %v", res.Error)
		return "", "", res.Error
	}

	// create basic app if not exsit
	if !appExist {
		appBasicInfo = appentity.ApplicationBasicInfo{
			ClusterId:        userInfo.ClusterId,
			GroupId:          userInfo.GroupId,
			UserId:           userInfo.UserId,
			AppName:          req.AppName,
			AppType:          constants.AppUserDefine,
			FunctionType:     req.FunctionType,
			LatestPubVersion: "",
			Description:      req.Description,
			CreateAt:         time.Now(),
			UpdateAt:         time.Now(),
		}

		err := appentity.CreateApplication(&appBasicInfo)
		if err != nil {
			logger.Errorf("CreateEdgeApplication create application failed. err = %v", err)
			return "", "", err
		}
	}

	// create app version if not exist
	data, err := json.Marshal(req.Network)
	if err != nil {
		logger.Errorf("CreateEdgeApplication Marshal network. err = %v", err)
		return "", "", err
	}
	netPolicy := string(data)

	archT := ""
	for _, v := range req.ArchType {
		archT = archT + v + ";"
	}
	archT = strings.TrimSuffix(archT, ";")
	if !verExist {
		appVersionInfo = appentity.ApplicationVersionInfo{
			AppName:               req.AppName,
			ClusterId:             userInfo.ClusterId,
			GroupId:               userInfo.GroupId,
			UserId:                userInfo.UserId,
			Version:               req.Version,
			Status:                appmodule.AppStatusUnpublished,
			ArchType:              archT,
			ContainerImage:        req.ContainerImage,
			ContainerImageVersion: req.ContainerImageVersion,
			ContainerImagePath:    imgPath,
			CpuQuota:              req.CpuQuota,
			MaxCpuQuota:           req.MaxCpuQuota,
			MemoryQuota:           req.MemoryQuota,
			MaxMemoryQuota:        req.MaxMemoryQuota,
			RestartPolicy:         req.RestartPolicy,
			Network:               netPolicy,
			CreateAt:              time.Now(),
			UpdateAt:              time.Now(),
			PublishAt:             time.Time{}.String(),
			OfflineAt:             time.Time{}.String(),
		}

		err := appentity.CreateApplicationVersion(&appVersionInfo)
		if err != nil {
			return "", "", err
		}
	}

	return appBasicInfo.AppName, appVersionInfo.Version, nil
}

// list edge application
func ListEdgeApplications(userInfo proto.ApulisHeader, req *appmodule.ListEdgeApplicationReq) (*[]appentity.ApplicationBasicInfo, int, error) {
	var appInfos []appentity.ApplicationBasicInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	var res *gorm.DB
	if req.AppType == appmodule.AppTypeAll {
		res = apulisdb.Db.Offset(offset).Limit(limit).
			Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
			Find(&appInfos)
	} else {
		res = apulisdb.Db.Offset(offset).Limit(limit).
			Where("ClusterId = ? and GroupId = ? and UserId = ? and AppType = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppType).
			Find(&appInfos)
	}

	if res.Error != nil {
		return &appInfos, total, res.Error
	}

	return &appInfos, int(res.RowsAffected), nil
}

// describe edge app
func DescribeEdgeApp(userInfo proto.ApulisHeader, req *appmodule.DescribeEdgeApplicationReq) (*appentity.ApplicationBasicInfo, error) {
	var appInfo appentity.ApplicationBasicInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		First(&appInfo)

	if res.Error != nil {
		return &appInfo, res.Error
	}

	return &appInfo, nil
}

// delete edge application
func DeleteEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.DeleteEdgeApplicationReq) error {
	var appInfo appentity.ApplicationBasicInfo

	// first: check if any deploy exist
	var total int64
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: check if any version exist
	apulisdb.Db.Model(&appentity.ApplicationVersionInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		Count(&total)
	if total != 0 {
		return appmodule.ErrApplicationVersionExist
	}

	// third: get application and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		First(&appInfo)
	if res.Error != nil {
		return res.Error
	}

	return appentity.DeleteApplication(&appInfo)
}
