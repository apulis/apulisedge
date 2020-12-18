// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	appservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

// list edge app deploys
func ListEdgeAppDeploys(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListEdgeAppDeployReq
	var appDeploys *[]appentity.ApplicationDeployInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node
	appDeploys, total, err = appservice.ListEdgeDeploys(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListEdgeAppDeployRsp{
		Total:      total,
		AppDeploys: appDeploys,
	}
	return SuccessResp(c, &req, data)
}

// deploy edge application
func DeployEdgeApp(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DeployEdgeApplicationReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// deploy application
	err = appservice.DeployEdgeApplication(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// undeploy edge application
func UnDeployEdgeApp(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.UnDeployEdgeApplicationReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// deploy application
	err = appservice.UnDeployEdgeApplication(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
