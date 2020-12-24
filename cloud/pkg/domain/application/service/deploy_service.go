// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"gorm.io/gorm"
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
	return listNodeCanDeployOrUpdate(userInfo, req.PageSize, req.PageNum, req.AppName, req.TargetVersion, "NodeName NOT IN ?")
}

func listNodeCanDeployOrUpdate(userInfo proto.ApulisHeader, pageSize, pageNum int, appName, targetVersion string, nodeQueryStr string) (*[]nodeentity.NodeBasicInfo, int, error) {
	var appVerInfo appentity.ApplicationVersionInfo
	var nodeInfos []nodeentity.NodeBasicInfo
	var deployInfos []appentity.ApplicationDeployInfo

	total := 0
	offset := pageSize * (pageNum - 1)
	limit := pageSize

	// get app version
	res := apulisdb.Db.Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
		userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, appName, targetVersion).
		First(&appVerInfo)
	if res.Error != nil {
		return &nodeInfos, total, res.Error
	}

	// get node already has deploy this version
	res = apulisdb.Db.Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
		userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, appName, targetVersion).
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
	var deployInfo appentity.ApplicationDeployInfo
	var newDeployInfo appentity.ApplicationDeployInfo
	var nodeInfo nodeentity.NodeBasicInfo

	// get cluster
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		logger.Infof("DeployEdgeApplication, can`t find cluster %d", userInfo.ClusterId)
		return cluster.ErrFindCluster
	}

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		logger.Infof("get application version info failed! err = %v", res.Error)
		return res.Error
	}

	archArr := strings.Split(appVerInfo.ArchType, ";")

	// check version status
	if appVerInfo.Status != appmodule.AppStatusPublished {
		return appmodule.ErrDeployStatusNotPublished
	}

	totalDeploy := 0
	var targetStatus string
	for _, node := range req.NodeNames {
		nodeInfo = nodeentity.NodeBasicInfo{}
		deployInfo = appentity.ApplicationDeployInfo{}

		// get node
		res = apulisdb.Db.
			Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ? and Arch IN ?",
				userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, node, archArr).
			First(&nodeInfo)
		if res.Error != nil {
			logger.Infof("get node failed! node = %s, err = %v", node, res.Error)
			continue
		}

		// get deploy
		res = apulisdb.Db.
			Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and NodeName = ?",
				userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, node).
			First(&deployInfo)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				targetStatus = constants.StatusInit
			} else {
				logger.Infof("get app deploy failed! node = %s, err = %v", node, err)
				continue
			}
		} else {
			targetStatus = constants.StatusUpdateInit
		}

		if targetStatus == constants.StatusInit {
			// store deploy info
			newDeploy := &appentity.ApplicationDeployInfo{
				ClusterId:          userInfo.ClusterId,
				GroupId:            userInfo.GroupId,
				UserId:             userInfo.UserId,
				AppName:            req.AppName,
				NodeName:           node,
				UniqueName:         nodeInfo.UniqueName,
				Version:            appVerInfo.Version,
				ContainerImagePath: appVerInfo.ContainerImagePath,
				Status:             constants.StatusInit,
				DeployUUID:         clu.GetUniqueName(cluster.ResourceDeploy),
				ContainerUUID:      clu.GetUniqueName(cluster.ResourceContainer),
				CreateAt:           time.Now(),
				UpdateAt:           time.Now(),
			}

			err = appentity.CreateAppDeploy(newDeploy)
			if err != nil {
				logger.Infof("create application deploy failed! node = %s, err = %v", node, err)
				continue
			}
			totalDeploy++
		} else if targetStatus == constants.StatusUpdateInit {
			if deployInfo.Version == req.Version {
				logger.Infof("deploy application deploy same version, skipped! node = %s, version = %s", node, deployInfo.Version)
				continue
			}

			if deployInfo.ContainerImagePath == appVerInfo.ContainerImagePath {
				logger.Infof("deploy application deploy same container, skipped! node = %s, version = %s, container = %s",
					node, deployInfo.Version, deployInfo.ContainerImagePath)
				continue
			}

			newDeployInfo = deployInfo
			newDeployInfo.Status = constants.StatusUpdateInit
			newDeployInfo.Version = req.Version
			newDeployInfo.ContainerImagePath = appVerInfo.ContainerImagePath
			newDeployInfo.UpdateAt = time.Now()

			err = appentity.UpdateAppDeploy(&newDeployInfo)
			if err != nil {
				logger.Infof("update application deploy failed! node = %s, err = %v", node, err)
				continue
			}
			totalDeploy++
		}
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
