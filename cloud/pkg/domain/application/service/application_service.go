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
