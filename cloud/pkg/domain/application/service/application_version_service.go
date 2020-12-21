// Copyright 2020 Apulis Technology Inc. All rights reserved.

package applicationservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"time"
)

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

// describe edge app version
func DescribeEdgeAppVersion(userInfo proto.ApulisHeader, req *appmodule.DescribeEdgeAppVersionReq) (*appentity.ApplicationVersionInfo, error) {
	var appVerInfo appentity.ApplicationVersionInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)

	if res.Error != nil {
		return &appVerInfo, res.Error
	}

	return &appVerInfo, nil
}

// publish app version
func PublishEdgeApplicationVersion(userInfo proto.ApulisHeader, req *appmodule.PublishEdgeApplicationVersionReq) error {
	var appInfo appentity.ApplicationBasicInfo
	var appVerInfo appentity.ApplicationVersionInfo

	// get app
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName).
		First(&appInfo)
	if res.Error != nil {
		return res.Error
	}

	// get app version
	res = apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	if appVerInfo.Status != appmodule.AppStatusUnpublished {
		return appmodule.ErrChangeAppVersionFailed
	}

	// app version update field: Status
	appVerInfo.Status = appmodule.AppStatusPublished
	appVerInfo.PublishAt = time.Now().String()
	err := appentity.UpdateApplicationVersion(&appVerInfo)
	if err != nil {
		return err
	}

	// app update field: LatestPubVersion
	appInfo.LatestPubVersion = req.Version
	return appentity.UpdateApplication(&appInfo)
}

// app version offline
func OfflineEdgeApplicationVersion(userInfo proto.ApulisHeader, req *appmodule.OfflineEdgeApplicationVersionReq) error {
	var appVerInfo appentity.ApplicationVersionInfo

	// get app version
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	if appVerInfo.Status != appmodule.AppStatusPublished {
		return appmodule.ErrChangeAppVersionFailed
	}

	appVerInfo.Status = appmodule.AppStatusOffline
	return appentity.UpdateApplicationVersion(&appVerInfo)
}

// delete edge application version
func DeleteEdgeApplicationVersion(userInfo proto.ApulisHeader, req *appmodule.DeleteEdgeApplicationVersionReq) error {
	var appVerInfo appentity.ApplicationVersionInfo

	// check if any deploy exist
	var total int64
	apulisdb.Db.Model(&appentity.ApplicationDeployInfo{}).
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		Count(&total)
	if total != 0 {
		return appmodule.ErrDeployExist
	}

	// get app version and delete
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and AppName = ? and Version = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.AppName, req.Version).
		First(&appVerInfo)
	if res.Error != nil {
		return res.Error
	}

	// check version status
	if appVerInfo.Status == appmodule.AppStatusPublished {
		return appmodule.ErrDeleteStatusPublished
	}

	return appentity.DeleteApplicationVersion(&appVerInfo)
}
