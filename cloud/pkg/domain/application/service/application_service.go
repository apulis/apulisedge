// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"fmt"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"time"
)

var logger = loggers.LogInstance()

// list edge application
func ListEdgeApplications(req *appmodule.ListEdgeApplicationReq) (*[]appentity.ApplicationBasicInfo, int, error) {
	var appInfos []appentity.ApplicationBasicInfo
	total := 0
	whereQueryStr := fmt.Sprintf("UserId = '%s'", req.UserId)
	res := apulisdb.Db.Offset(req.PageNum).Limit(req.PageSize).Where(whereQueryStr).Find(&appInfos)

	if res.Error != nil {
		return &appInfos, total, res.Error
	}

	return &appInfos, int(res.RowsAffected), nil
}

// create edge application
func CreateEdgeApplication(req *appmodule.CreateEdgeApplicationReq) (*appentity.ApplicationBasicInfo, error) {
	appInfo := &appentity.ApplicationBasicInfo{
		UserId:                req.UserId,
		UserName:              req.UserName,
		AppName:               req.AppName,
		ArchType:              req.ArchType,
		Version:               req.Version,
		ContainerImage:        req.ContainerImage,
		ContainerImageVersion: req.ContainerImageVersion,
		ContainerImagePath:    req.ContainerImagePath,
		CpuQuota:              req.CpuQuota,
		MemoryQuota:           req.MemoryQuota,
		CreateAt:              time.Now(),
		UpdateAt:              time.Now(),
	}

	return appInfo, appentity.CreateApplication(appInfo)
}

// delete edge application
func DeleteEdgeApplication(req *appmodule.DeleteEdgeApplicationReq) error {
	// first: check if any deploy exist
	var total int64
	whereQueryStr := fmt.Sprintf("UserId = '%s' and AppName = '%s' and Version = '%s'", req.UserId, req.AppName, req.Version)
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).Where(whereQueryStr).Count(&total)

	if total != 0 {
		return appmodule.ErrDeployExist
	}

	appInfo, err := appentity.GetApplication(req.UserId, req.AppName, req.Version)
	if err != nil {
		return err
	}

	return appentity.DeleteApplication(appInfo)
}

// list edge deploys
func ListEdgeDeploys(req *appmodule.ListEdgeAppDeployReq) (*[]appentity.ApplicationDeployInfo, int, error) {
	var appDeloys []appentity.ApplicationDeployInfo
	total := 0
	whereQueryStr := fmt.Sprintf("UserId = '%s' and AppName = '%s' and Version = '%s'", req.UserId, req.AppName, req.Version)
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
	var appInfo appentity.ApplicationBasicInfo

	whereQueryStr := fmt.Sprintf("UserId = '%s' and AppName = '%s'", req.UserId, req.AppName)
	res := apulisdb.Db.Where(whereQueryStr).First(&appInfo)
	if res.Error != nil {
		return res.Error
	}

	// store deploy info
	hostPort := 0
	if req.PortMapping.Enable {
		hostPort = req.PortMapping.HostPort
	}

	deployInfo := &appentity.ApplicationDeployInfo{
		UserId:                req.UserId,
		UserName:              req.UserName,
		AppName:               req.AppName,
		NodeName:              req.NodeName,
		ArchType:              appInfo.ArchType,
		Version:               appInfo.Version,
		Status:                constants.StatusInit,
		Namespace:             req.NamespaceName,
		ContainerImage:        appInfo.ContainerImage,
		ContainerImageVersion: appInfo.ContainerImageVersion,
		ContainerImagePath:    appInfo.ContainerImagePath,
		CpuQuota:              appInfo.CpuQuota,
		MemoryQuota:           appInfo.MemoryQuota,
		ContainerPort:         req.ContainerPort,
		HostPort:              hostPort,
		CreateAt:              time.Now(),
		UpdateAt:              time.Now(),
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

	whereQueryStr := fmt.Sprintf("UserId = '%s' and AppName = '%s' and NodeName = '%s' and Version = '%s'", req.UserId, req.AppName, req.NodeName, req.Version)
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
