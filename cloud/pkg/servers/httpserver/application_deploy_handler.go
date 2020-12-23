// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	appmodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/application"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	appservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/service"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
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

	// list app deploy
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

// list node deploys
func ListNodeDeploys(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListNodeDeployReq
	var appDeploys *[]appentity.ApplicationDeployInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node deploy
	appDeploys, total, err = appservice.ListNodeDeploys(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListNodeDeployRsp{
		Total:      total,
		AppDeploys: appDeploys,
	}
	return SuccessResp(c, &req, data)
}

// list node can deploy
func ListNodeCanDeploy(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListNodeCanDeployReq
	var nodes *[]nodeentity.NodeBasicInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node deploy
	nodes, total, err = appservice.ListNodeCanDeploy(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListNodeCanDeployRsp{
		Total: total,
		Nodes: nodes,
	}
	return SuccessResp(c, &req, data)
}

// list node can update
func ListNodeCanUpdate(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListNodeCanUpdateReq
	var nodes *[]nodeentity.NodeBasicInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node update
	nodes, total, err = appservice.ListNodeCanUpdate(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListNodeCanUpdateRsp{
		Total: total,
		Nodes: nodes,
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

// update deploy edge application
func UpdateDeployEdgeApp(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.UpdateDeployEdgeApplicationReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// deploy application
	err = appservice.UpdateDeployEdgeApplication(*userInfo, &reqContent)
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
