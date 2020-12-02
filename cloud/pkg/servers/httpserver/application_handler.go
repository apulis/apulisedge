// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	appservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ApplicationHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/application")

	group.POST("/createNode", wrapper(CreateEdgeApplication))

}

// create edge application
func CreateEdgeApplication(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.CreateEdgeApplicationReq
	var app *appentity.ApplicationBasicInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// create node
	app, err = appservice.CreateEdgeApplication(&reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.CreateEdgeApplicationRsp{
		Application: app,
	}
	return SuccessResp(c, &req, data)
}
