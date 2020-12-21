// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

func ListContainerImageVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.ListContainerImageVersionReq
	var imageVers *[]imageentity.UserContainerImageVersionInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list image version
	imageVers, total, err = imageservice.ListContainerImageVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.ListContainerImageVersionRsp{
		Total:         total,
		ImageVersions: imageVers,
	}

	return SuccessResp(c, &req, data)
}

// describe image version
func DescribeContainerImageVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DescribeContainerImageVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// describe application
	imgVer, err := imageservice.DescribeContainerImageVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.DescribeContainerImageVersionRsp{
		ImageVersion: imgVer,
	}

	return SuccessResp(c, &req, data)
}

// delete image version
func DeleteImageVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DeleteContainerImageVersionReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// delete image version
	err = imageservice.DeleteContainterImageVersion(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
