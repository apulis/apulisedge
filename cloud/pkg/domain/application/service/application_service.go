// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"fmt"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"gorm.io/gorm"
	"time"
)

var logger = loggers.LogInstance()

// create edge application
// this interface can both create basic app and app version
func CreateEdgeApplication(req *appmodule.CreateEdgeApplicationReq) (string, string, error) {
	appBasicInfo, err := appentity.GetApplication(req.ClusterId, req.GroupId, req.UserId, req.AppName)
	if err == gorm.ErrRecordNotFound {
		// create basic app if not exsit
		appBasicInfo = &appentity.ApplicationBasicInfo{
			ClusterId:        req.ClusterId,
			GroupId:          req.GroupId,
			UserId:           req.UserId,
			AppName:          req.AppName,
			AppType:          constants.AppUserDefine,
			FunctionType:     req.FunctionType,
			LatestPubVersion: "",
			Description:      req.Description,
			CreateAt:         time.Now(),
			UpdateAt:         time.Now(),
		}

		err := appentity.CreateApplication(appBasicInfo)
		if err != nil {
			logger.Errorf("CreateEdgeApplication CreateApplication failed. err = %v", err)
			return "", "", err
		}
	} else if err != nil {
		logger.Errorf("CreateEdgeApplication GetApplication failed. err = %v", err)
		return "", "", err
	}

	// create app version
	logger.Infof("%v", req)
	appVersionInfo := &appentity.ApplicationVersionInfo{
		AppName:               appBasicInfo.AppName,
		ClusterId:             req.ClusterId,
		GroupId:               req.GroupId,
		UserId:                req.UserId,
		Version:               req.Version,
		ArchType:              req.ArchType,
		ContainerImage:        req.ContainerImage,
		ContainerImageVersion: req.ContainerImageVersion,
		ContainerImagePath:    req.ContainerImagePath,
		CpuQuota:              req.CpuQuota,
		MaxCpuQuota:           req.MaxCpuQuota,
		MemoryQuota:           req.MemoryQuota,
		MaxMemoryQuota:        req.MaxMemoryQuota,
		CreateAt:              time.Now(),
		UpdateAt:              time.Now(),
	}

	err = appentity.CreateApplicationVersion(appVersionInfo)
	if err != nil {
		return "", "", err
	}

	return appBasicInfo.AppName, appVersionInfo.Version, nil
}

// list edge application
func ListEdgeApplications(req *appmodule.ListEdgeApplicationReq) (*[]appentity.ApplicationBasicInfo, int, error) {
	var appInfos []appentity.ApplicationBasicInfo
	total := 0
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s'", req.ClusterId, req.GroupId, req.UserId)
	res := apulisdb.Db.Offset(req.PageNum).Limit(req.PageSize).Where(whereQueryStr).Find(&appInfos)

	if res.Error != nil {
		return &appInfos, total, res.Error
	}

	return &appInfos, int(res.RowsAffected), nil
}

// delete edge application
func DeleteEdgeApplication(req *appmodule.DeleteEdgeApplicationReq) error {
	// first: check if any deploy exist
	var total int64
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName)
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).Where(whereQueryStr).Count(&total)

	if total != 0 {
		return appmodule.ErrDeployExist
	}

	appInfo, err := appentity.GetApplication(req.ClusterId, req.GroupId, req.UserId, req.AppName)
	if err != nil {
		return err
	}

	return appentity.DeleteApplication(appInfo)
}

// list edge application versions
func ListEdgeApplicationVersions(req *appmodule.ListEdgeApplicationVersionReq) (*[]appentity.ApplicationVersionInfo, int, error) {
	var appVerInfos []appentity.ApplicationVersionInfo
	total := 0
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName)
	res := apulisdb.Db.Offset(req.PageNum).Limit(req.PageSize).Where(whereQueryStr).Find(&appVerInfos)

	if res.Error != nil {
		return &appVerInfos, total, res.Error
	}

	return &appVerInfos, int(res.RowsAffected), nil
}

// delete edge application version
func DeleteEdgeApplicationVersion(req *appmodule.DeleteEdgeApplicationVersionReq) error {
	// first: check if any deploy exist
	var total int64
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s' and Version = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName, req.Version)
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).Where(whereQueryStr).Count(&total)

	if total != 0 {
		return appmodule.ErrDeployExist
	}

	appVerInfo, err := appentity.GetApplicationVersion(req.ClusterId, req.GroupId, req.UserId, req.AppName, req.Version)
	if err != nil {
		return err
	}

	return appentity.DeleteApplicationVersion(appVerInfo)
}

// list edge deploys
func ListEdgeDeploys(req *appmodule.ListEdgeAppDeployReq) (*[]appentity.ApplicationDeployInfo, int, error) {
	var appDeloys []appentity.ApplicationDeployInfo
	total := 0
	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s' and Version = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName, req.Version)
	res := apulisdb.Db.Offset(req.PageNum).Limit(req.PageSize).Where(whereQueryStr).Find(&appDeloys)

	if res.Error != nil {
		return &appDeloys, total, res.Error
	}

	return &appDeloys, int(res.RowsAffected), nil
}

// deploy edge application
func DeployEdgeApplication(req *appmodule.DeployEdgeApplicationReq) error {
	// get application
	var err error
	var appVerInfo appentity.ApplicationVersionInfo

	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s' and Version = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName, req.Version)
	res := apulisdb.Db.Where(whereQueryStr).First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	// store deploy info
	deployInfo := &appentity.ApplicationDeployInfo{
		ClusterId: req.ClusterId,
		GroupId:   req.GroupId,
		UserId:    req.UserId,
		AppName:   req.AppName,
		NodeName:  req.NodeName,
		Version:   appVerInfo.Version,
		Status:    constants.StatusInit,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}

	err = appentity.CreateAppDeploy(deployInfo)
	if err != nil {
		logger.Infof("create application deploy failed! err = %v", err)
		return err
	}

	return nil
}

// undeploy edge application
func UnDeployEdgeApplication(req *appmodule.UnDeployEdgeApplicationReq) error {
	// get application
	var err error
	var appDeployInfo appentity.ApplicationDeployInfo

	whereQueryStr := fmt.Sprintf("ClusterId = '%s' and GroupId = '%s' and UserId = '%s' and AppName = '%s' and NodeName = '%s' and Version = '%s'",
		req.ClusterId, req.GroupId, req.UserId, req.AppName, req.NodeName, req.Version)
	res := apulisdb.Db.Where(whereQueryStr).First(&appDeployInfo)
	if res.Error != nil {
		return res.Error
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
