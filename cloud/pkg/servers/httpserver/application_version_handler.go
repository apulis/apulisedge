// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	appservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

// list edge app versions
func ListEdgeAppVersions(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListEdgeApplicationVersionReq
	var appVers *[]appentity.ApplicationVersionInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node
	appVers, total, err = appservice.ListEdgeApplicationVersions(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListEdgeApplicationVersionRsp{
		Total:       total,
		AppVersions: appVers,
	}
	return SuccessResp(c, &req, data)
}

// describe edge application version
func DescribeEdgeAppVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DescribeEdgeAppVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// describe application version
	appVer, err := appservice.DescribeEdgeAppVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.DescribeEdgeAppVersionRsp{
		AppVersion: appVer,
	}

	return SuccessResp(c, &req, data)
}

// publish app version
func PublishEdgeAppVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.PublishEdgeApplicationVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// publish application
	err = appservice.PublishEdgeApplicationVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// offline app version
func OfflineEdgeAppVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.OfflineEdgeApplicationVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// publish application
	err = appservice.OfflineEdgeApplicationVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// delete edge application version
func DeleteEdgeAppVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DeleteEdgeApplicationVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// delete application
	err = appservice.DeleteEdgeApplicationVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
