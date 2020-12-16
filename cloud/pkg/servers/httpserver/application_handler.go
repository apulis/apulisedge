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

	// add authentication
	group.Use(Auth())

	/// edge application
	group.POST("/listApplication", wrapper(ListEdgeApps))
	group.POST("/createApplication", wrapper(CreateEdgeApplication))
	group.POST("/deleteApplication", wrapper(DeleteEdgeApplication))

	// edge application version
	group.POST("/listApplicationVersion", wrapper(ListEdgeAppVersions))
	group.POST("/deleteApplicationVersion", wrapper(DeleteEdgeApplicationVersion))

	// application deployment
	group.POST("/listApplicationDeploy", wrapper(ListEdgeAppDeploys))
	group.POST("/deployApplication", wrapper(DeployEdgeApplication))
	group.POST("/undeployApplication", wrapper(UnDeployEdgeApplication))
}

// list edge apps
func ListEdgeApps(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListEdgeApplicationReq
	var apps *[]appentity.ApplicationBasicInfo
	var total int

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// list node
	apps, total, err = appservice.ListEdgeApplications(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListEdgeApplicationRsp{
		Total: total,
		Apps:  apps,
	}
	return SuccessResp(c, &req, data)
}

// create edge application
// this interface can both create basic app and app version
func CreateEdgeApplication(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.CreateEdgeApplicationReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// create application
	appCreated, verCreated, err := appservice.CreateEdgeApplication(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.CreateEdgeApplicationRsp{
		AppCreated:     appCreated,
		VersionCreated: verCreated,
	}

	return SuccessResp(c, &req, data)
}

// delete edge application
func DeleteEdgeApplication(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DeleteEdgeApplicationReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// delete application
	err = appservice.DeleteEdgeApplication(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// list edge app versions
func ListEdgeAppVersions(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListEdgeApplicationVersionReq
	var appVers *[]appentity.ApplicationVersionInfo
	var total int

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// list node
	appVers, total, err = appservice.ListEdgeApplicationVersions(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := appmodule.ListEdgeApplicationVersionRsp{
		Total:       total,
		AppVersions: appVers,
	}
	return SuccessResp(c, &req, data)
}

// delete edge application version
func DeleteEdgeApplicationVersion(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DeleteEdgeApplicationVersionReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// delete application
	err = appservice.DeleteEdgeApplicationVersion(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// list edge app deploys
func ListEdgeAppDeploys(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.ListEdgeAppDeployReq
	var appDeploys *[]appentity.ApplicationDeployInfo
	var total int

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// list node
	appDeploys, total, err = appservice.ListEdgeDeploys(userInfo, &reqContent)
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
func DeployEdgeApplication(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.DeployEdgeApplicationReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// deploy application
	err = appservice.DeployEdgeApplication(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// undeploy edge application
func UnDeployEdgeApplication(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent appmodule.UnDeployEdgeApplicationReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	// deploy application
	err = appservice.UnDeployEdgeApplication(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
