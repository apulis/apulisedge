// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"context"
	"fmt"
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ImageHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/image")

	// add authentication
	group.Use(Auth())

	// image
	group.POST("/listImage", wrapper(ListContainerImage))
	group.POST("/uploadImage", wrapper(UploadContainerImage))
	group.POST("/deleteImage", wrapper(DeleteImage))

	// image version
	group.POST("/listImageVersion", wrapper(ListContainerImageVersion))
	group.POST("/deleteImageVersion", wrapper(DeleteImageVersion))
}

func ListContainerImage(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.ListContainerImageReq
	var images []imageentity.UserContainerImageInfo
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

	// list image
	images, total, err = imageservice.ListContainerImage(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.ListContainerImageRsp{
		Total:  total,
		Images: []imagemodule.RspContainerImageInfo{},
	}

	for i := 0; i < data.Total; i++ {
		img := imagemodule.RspContainerImageInfo{
			ClusterId:    images[i].ClusterId,
			GroupId:      images[i].GroupId,
			UserId:       images[i].UserId,
			ImageName:    images[i].ImageName,
			OrgName:      images[i].OrgName,
			VersionCount: 1,
			UpdateAt:     images[i].UpdateAt,
		}
		data.Images = append(data.Images, img)
	}

	return SuccessResp(c, &req, data)
}

// upload container image
func UploadContainerImage(c *gin.Context) error {
	var err error

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	// org
	orgName := c.PostForm("orgName")
	if orgName == "" {
		return NoReqAppError(c, ErrOrgNameNeeded.Error())
	}

	// single file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return NoReqAppError(c, err.Error())
	}
	logger.Infof("Uploading container image, file = %s", fileHeader.Filename)

	dstFile := "/tmp/apulis/images/" + fileHeader.Filename
	err = c.SaveUploadedFile(fileHeader, dstFile)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	// get cluster
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		logger.Infof("UploadContainerImage, can`t find cluster %d", userInfo.ClusterId)
		return NoReqAppError(c, err.Error())
	}

	// image load
	ctx := context.Background()
	cli, err := clu.NewDockerClient()
	if err != nil {
		return NoReqAppError(c, err.Error())
	}
	defer clu.CloseDockerClient(cli)

	img, err := clu.DockerImageLoad(ctx, cli, dstFile)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}
	logger.Infof("Image load succ, load tag = %s", img)

	// image tag
	tag, ver, err := clu.GetImageNameAndVersion(img)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	dstImage := clu.GetHarborAddress() + "/" + clu.GetHarborProject() + "/" + orgName + "/" + tag + ":" + ver
	err = clu.DockerImageTag(ctx, cli, tag+":"+ver, dstImage)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	// img push
	err = clu.DockerImagePush(ctx, cli, dstImage)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}
	logger.Infof("Image push succ, tag = %s", img)

	// add to db
	err = imageservice.AddContainerImage(userInfo, orgName, tag, ver, dstImage)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	return NoReqSuccessResp(c, fmt.Sprintf("'%s' uploaded!", fileHeader.Filename))
}

// delete image
func DeleteImage(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DeleteContainerImageReq

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

	// delete image
	err = imageservice.DeleteContainerImage(userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

// image version =======================================================================================
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
