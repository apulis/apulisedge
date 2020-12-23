// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/satori/go.uuid"
	"strings"
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

// list node can deploy
func ListNodeCanDeploy(userInfo proto.ApulisHeader, req *appmodule.ListNodeCanDeployReq) (*[]nodeentity.NodeBasicInfo, int, error) {
	return listNodeCanDeployOrUpdate(userInfo, req.PageSize, req.PageNum, req.AppName, req.Version, "NodeName NOT IN ?")
}

// list node can update
func ListNodeCanUpdate(userInfo proto.ApulisHeader, req *appmodule.ListNodeCanUpdateReq) (*[]nodeentity.NodeBasicInfo, int, error) {
	return listNodeCanDeployOrUpdate(userInfo, req.PageSize, req.PageNum, req.AppName, req.Version, "NodeName IN ?")
}

func listNodeCanDeployOrUpdate(userInfo proto.ApulisHeader, pageSize, pageNum int, appName, version string, nodeQueryStr string) (*[]nodeentity.NodeBasicInfo, int, error) {
	var appVerInfo appentity.ApplicationVersionInfo
	var nodeInfos []nodeentity.NodeBasicInfo
	var deployInfos []appentity.ApplicationDeployInfo

	total := 0
	offset := pageSize * (pageNum - 1)
	limit := pageSize

	// get app version
	res := apulisdb.Db.Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
		userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, appName, version).
		First(&appVerInfo)
	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	// get node already has deploy
	res = apulisdb.Db.Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?",
		userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, appName).
		Find(&deployInfos)
	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	var nodeArr []string
	for _, v := range deployInfos {
		nodeArr = append(nodeArr, v.NodeName)
	}

	// add a empty string, otherwise the (NOT IN)/(IN) statement will not work
	if len(nodeArr) == 0 {
		nodeArr = append(nodeArr, "")
	}
	archArr := strings.Split(appVerInfo.ArchType, ";")

	// get node which arch is equal to app version
	res = apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and Status = ? and Arch IN ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, nodemodule.StatusOnline, archArr, nodeArr).
		Where(nodeQueryStr, nodeArr).
		Find(&nodeInfos)
	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	return &nodeInfos, total, res.Error
}

// deploy edge application
func DeployEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.DeployEdgeApplicationReq) error {
	// get application version
	var err error
	var appVerInfo appentity.ApplicationVersionInfo
	var nodeInfo nodeentity.NodeBasicInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		logger.Infof("get application version info failed! err = %v", res.Error)
		return res.Error
	}

	// check version status
	if appVerInfo.Status != appmodule.AppStatusPublished {
		return appmodule.ErrDeployStatusNotPublished
	}

	totalDeploy := 0
	for _, node := range req.NodeNames {
		res = apulisdb.Db.
			Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?",
				userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, node).
			First(&nodeInfo)
		if res.Error != nil {
			logger.Infof("get node failed! node = %s, err = %v", node, res.Error)
			continue
		}

		// store deploy info
		deployInfo := &appentity.ApplicationDeployInfo{
			ClusterId:  userInfo.ClusterId,
			GroupId:    userInfo.GroupId,
			UserId:     userInfo.UserId,
			AppName:    req.AppName,
			NodeName:   node,
			UniqueName: nodeInfo.UniqueName,
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

// update deploy edge application
func UpdateDeployEdgeApplication(userInfo proto.ApulisHeader, req *appmodule.UpdateDeployEdgeApplicationReq) error {
	// get application version
	var err error
	var appVerInfo appentity.ApplicationVersionInfo
	var deployInfo appentity.ApplicationDeployInfo
	var newDeployInfo appentity.ApplicationDeployInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.TargetVersion).
		First(&appVerInfo)
	if res.Error != nil {
		logger.Infof("get application version info failed! err = %v", res.Error)
		return res.Error
	}

	// check version status
	if appVerInfo.Status != appmodule.AppStatusPublished {
		return appmodule.ErrDeployStatusNotPublished
	}

	totalUpdate := 0
	for _, node := range req.NodeNames {
		// get deploy
		res = apulisdb.Db.
			Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and NodeName = ?",
				userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, node).
			First(&deployInfo)
		if res.Error != nil {
			logger.Infof("update application deploy failed! node = %s, err = %v", node, err)
			continue
		}

		if deployInfo.Version == req.TargetVersion {
			logger.Infof("update application deploy same version! node = %s, version = %s", node, deployInfo.Version)
			continue
		}

		newDeployInfo = deployInfo
		newDeployInfo.Status = constants.StatusUpdating
		newDeployInfo.Version = req.TargetVersion
		newDeployInfo.UpdateAt = time.Now()

		err = appentity.UpdateAppDeploy(&newDeployInfo)
		if err != nil {
			logger.Infof("update application deploy failed! node = %s, err = %v", node, err)
			continue
		}
		totalUpdate++
	}

	if totalUpdate == len(req.NodeNames) {
		return nil
	} else if totalUpdate > 0 {
		return appmodule.ErrUpdatePartFails
	} else {
		return appmodule.ErrUpdateAllFails
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
