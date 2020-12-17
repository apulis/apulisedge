// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ListContainerImageVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.ListContainerImageVersionReq
	var imageVers *[]imageentity.UserContainerImageVersionInfo
	var total int

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// validate request content
	validate := validator.New()
	err = validate.Struct(reqContent)
	if err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// list image version
	imageVers, total, err = imageservice.ListContainerImageVersion(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.ListContainerImageVersionRsp{
		Total:         total,
		ImageVersions: imageVers,
	}

	return SuccessResp(c, &req, data)
}

// delete image version
func DeleteImageVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DeleteContainerImageVersionReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// validate request content
	validate := validator.New()
	err = validate.Struct(reqContent)
	if err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// delete image version
	err = imageservice.DeleteContainterImageVersion(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
