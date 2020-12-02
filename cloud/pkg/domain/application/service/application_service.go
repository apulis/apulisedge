// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
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
