// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// pre handler is used to handle the common logic of each comming request
func PreHandler(c *gin.Context, req *proto.Message, reqContent interface{}) (*proto.ApulisHeader, *APIErrorResp) {
	var err error

	// check reqContent type
	contentType := reflect.TypeOf(reqContent)
	switch contentType.Kind() {
	case reflect.Ptr:
		break
	default:
		logger.Errorf("PreHandler: reqContent is not a pointer")
		return nil, ServerError(c, req)
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return nil, ParameterError(c, req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), reqContent); err != nil {
		return nil, ParameterError(c, req, err.Error())
	}

	// validate request content
	validate := validator.New()
	err = validate.Struct(reqContent)
	if err != nil {
		return nil, ParameterError(c, req, err.Error())
	}

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return nil, AppError(c, req, APP_ERROR_CODE, err.Error())
	}

	return &userInfo, nil
}
