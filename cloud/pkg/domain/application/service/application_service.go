// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	// check if i have this image
	imgPath, imgExsit := imageservice.DoIHaveTheImageVersion(userInfo, req.OrgName, req.ContainerImage, req.ContainerImageVersion)
	if !imgExsit {
		return "", "", appmodule.ErrImageVersionNotExist
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
	if !verExist {
		appVersionInfo = appentity.ApplicationVersionInfo{
			AppName:               req.AppName,
			ClusterId:             userInfo.ClusterId,
			GroupId:               userInfo.GroupId,
			UserId:                userInfo.UserId,
			Version:               req.Version,
			ArchType:              req.ArchType,
			ContainerImage:        req.ContainerImage,
			ContainerImageVersion: req.ContainerImageVersion,
			ContainerImagePath:    imgPath,
			CpuQuota:              req.CpuQuota,
			MaxCpuQuota:           req.MaxCpuQuota,
			MemoryQuota:           req.MemoryQuota,
			MaxMemoryQuota:        req.MaxMemoryQuota,
			CreateAt:              time.Now(),
			UpdateAt:              time.Now(),
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

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Find(&appInfos)

	if res.Error != nil {
		return &appInfos, total, res.Error
	}

	return &appInfos, int(res.RowsAffected), nil
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

// list edge application versions
func ListEdgeApplicationVersions(userInfo proto.ApulisHeader, req *appmodule.ListEdgeApplicationVersionReq) (*[]appentity.ApplicationVersionInfo, int, error) {
	var appVerInfos []appentity.ApplicationVersionInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		Find(&appVerInfos)

	if res.Error != nil {
		return &appVerInfos, total, res.Error
	}

	return &appVerInfos, int(res.RowsAffected), nil
}

// delete edge application version
func DeleteEdgeApplicationVersion(userInfo proto.ApulisHeader, req *appmodule.DeleteEdgeApplicationVersionReq) error {
	var appVerInfo appentity.ApplicationVersionInfo

	// first: check if any deploy exist
	var total int64
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// second: get app version and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	return appentity.DeleteApplicationVersion(&appVerInfo)
}

// list edge deploys
func ListEdgeDeploys(userInfo proto.ApulisHeader, req *appmodule.ListEdgeAppDeployReq) (*[]appentity.ApplicationDeployInfo, int, error) {
	var appDeloys []appentity.ApplicationDeployInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		Find(&appDeloys)
	if res.Error != nil {
		return &appDeloys, total, res.Error
	}

	return &appDeloys, int(res.RowsAffected), nil
}

// deploy edge application
func DeployEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.DeployEdgeApplicationReq) error {
	// get application version
	var err error
	var appVerInfo appentity.ApplicationVersionInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		logger.Infof("create application deploy failed! err = %v", res.Error)
		return res.Error
	}

	// store deploy info
	deployInfo := &appentity.ApplicationDeployInfo{
		ClusterId:  userInfo.ClusterId,
		GroupId:    userInfo.GroupId,
		UserId:     userInfo.UserId,
		AppName:    req.AppName,
		NodeName:   req.NodeName,
		Version:    appVerInfo.Version,
		Status:     constants.StatusInit,
		DeployUUID: uuid.NewV4().String(),
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	}

	err = appentity.CreateAppDeploy(deployInfo)
	if err != nil {
		logger.Infof("create application deploy failed! err = %v", err)
		return err
	}

	return nil
}

// undeploy edge application
func UnDeployEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.UnDeployEdgeApplicationReq) error {
	// get application
	var err error
	var appDeployInfo appentity.ApplicationDeployInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and NodeName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.NodeName, req.Version).
		First(&appDeployInfo)
	if res.Error != nil {
		return res.Error
	}

	if appDeployInfo.Status == constants.StatusDeleting {
		return constants.ErrUnDeploying
	}

	// modify status directly
	appDeployInfo.Status = constants.StatusDeleting

	err = appentity.UpdateAppDeploy(&appDeployInfo)
	if err != nil {
		logger.Infof("delete application deploy failed! err = %v", err)
		return err
	}

	return nil
}
