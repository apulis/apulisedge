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

// list node deploys
func ListNodeDeploys(userInfo proto.ApulisHeader, req *appmodule.ListNodeDeployReq) (*[]appentity.ApplicationDeployInfo, int, error) {
	var appDeloys []appentity.ApplicationDeployInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
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

	totalDeploy := 0
	for _, node := range req.NodeNames {
		// store deploy info
		deployInfo := &appentity.ApplicationDeployInfo{
			ClusterId:  userInfo.ClusterId,
			GroupId:    userInfo.GroupId,
			UserId:     userInfo.UserId,
			AppName:    req.AppName,
			NodeName:   node,
			Version:    appVerInfo.Version,
			Status:     constants.StatusInit,
			DeployUUID: uuid.NewV4().String(),
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
		}

		err = appentity.CreateAppDeploy(deployInfo)
		if err != nil {
			logger.Infof("create application deploy failed! node = %s, err = %v", node, err)
			continue
		}
		totalDeploy++
	}

	if totalDeploy == len(req.NodeNames) {
		return nil
	} else if totalDeploy > 0 {
		return appmodule.ErrDeployPartFails
	} else {
		return appmodule.ErrDeployAllFails
	}
}

// undeploy edge application
func UnDeployEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.UnDeployEdgeApplicationReq) error {
	// get application deploy
	var err error
	var appDeployInfos []appentity.ApplicationDeployInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and NodeName IN ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.NodeNames, req.Version).
		Find(&appDeployInfos)
	if res.Error != nil {
		return res.Error
	}

	totalUnDeploy := 0
	for _, info := range appDeployInfos {
		if info.Status == constants.StatusDeleting {
			continue
		}

		info.Status = constants.StatusDeleting
		err = appentity.UpdateAppDeploy(&info)
		if err != nil {
			logger.Infof("delete application deploy failed! node = %s, err = %v", info.NodeName, err)
			continue
		}

		totalUnDeploy++
	}

	if totalUnDeploy == len(req.NodeNames) {
		return nil
	} else if totalUnDeploy > 0 {
		return appmodule.ErrUnDeployPartFails
	} else {
		return appmodule.ErrUnDeployAllFails
	}
}
