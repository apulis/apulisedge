// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"fmt"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	imageservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ImageHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/image")

	// add authentication
	//group.Use(Auth())

	group.POST("/listImage", wrapper(ListContainerImage))
	group.POST("/uploadImage", wrapper(UploadContainerImage))
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

	// TODO validate reqContent

	// list node
	images, total, err = imageservice.ListContainerImage(&reqContent)
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

	// single file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	logger.Infof("Upload container image, file = %s", fileHeader.Filename)

	err = c.SaveUploadedFile(fileHeader, "/tmp/apulis/images/"+fileHeader.Filename)
	if err != nil {
		return NoReqAppError(c, err.Error())
	}

	return NoReqSuccessResp(c, fmt.Sprintf("'%s' uploaded!", fileHeader.Filename))
}
