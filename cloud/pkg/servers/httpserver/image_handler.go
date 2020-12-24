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
	"github.com/gin-gonic/gin/binding"
)

func ImageHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/image")

	// add authentication
	group.Use(Auth())

	// image
	group.POST("/listImage", wrapper(ListContainerImage))
	group.POST("/describeImage", wrapper(DescribeContainerImage))
	group.POST("/uploadImage", wrapper(UploadContainerImage))
	group.POST("/deleteImage", wrapper(DeleteImage))

	// image version
	group.POST("/listImageVersion", wrapper(ListContainerImageVersion))
	group.POST("/describeImageVersion", wrapper(DescribeContainerImageVersion))
	group.POST("/deleteImageVersion", wrapper(DeleteImageVersion))

	// org
	group.POST("/createOrg", wrapper(CreateImageOrg))
	group.POST("/listOrg", wrapper(ListImageOrg))
	group.POST("/deleteOrg", wrapper(DeleteImageOrg))
}

func ListContainerImage(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.ListContainerImageReq
	var images []imageentity.UserContainerImageInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list image
	images, total, err = imageservice.ListContainerImage(*userInfo, &reqContent)
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

// describe image
func DescribeContainerImage(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent imagemodule.DescribeContainerImageReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// describe application
	img, err := imageservice.DescribeContainerImage(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := imagemodule.DescribeContainerImageRsp{
		Image: img,
	}

	return SuccessResp(c, &req, data)
}

// upload container image
func UploadContainerImage(c *gin.Context) error {
	var err error
	var reqContent imagemodule.UploadContainerImageReq

	if err = c.ShouldBindWith(&reqContent, binding.FormMultipart); err != nil {
		return NoReqAppError(c, err.Error())
	}

	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	if !imageservice.DoIHaveTheOrg(userInfo, reqContent.OrgName) {
		return NoReqAppError(c, ErrIDontHaveOrg.Error())
	}

	fileHeader := reqContent.File
	orgName := reqContent.OrgName

	// single file
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

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// delete image
	err = imageservice.DeleteContainerImage(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
