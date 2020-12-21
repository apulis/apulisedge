// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/satori/go.uuid"
	"time"
)

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

	// check version status
	if appVerInfo.Status != appmodule.AppStatusPublished {
		return appmodule.ErrDeployStatusNotPublished
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
