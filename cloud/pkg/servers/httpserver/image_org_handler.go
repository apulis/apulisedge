// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

// create org
func CreateImageOrg(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.CreateContainerImageOrgReq
	var org *imageentity.ContainerImageOrg

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// create node
	org, err = imageservice.CreateContainerImageOrg(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.CreateContainerImageOrgRsp{
		Org: org,
	}
	return SuccessResp(c, &req, data)
}

// list org
func ListImageOrg(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.ListContainerImageOrgReq
	var orgs *[]imageentity.ContainerImageOrg
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	orgs, total, err = imageservice.ListContainerImageOrg(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.ListContainerImageOrgRsp{
		Total:     total,
		ImageOrgs: orgs,
	}
	return SuccessResp(c, &req, data)
}

// delete org
func DeleteImageOrg(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DeleteContainerImageOrgReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// delete org
	err = imageservice.DeleteContainterImageOrg(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
