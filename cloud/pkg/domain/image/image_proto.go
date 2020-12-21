// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"mime/multipart"
	"time"
)

// TODO add param validate, like node.ListEdgeNodesReq

// List image
type ListContainerImageReq struct {
	PageNum  int `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int `json:"pageSize" validate:"gte=1,lte=1000"`
}

type RspContainerImageInfo struct {
	ClusterId    int64     `json:"clusterId"`
	GroupId      int64     `json:"groupId"`
	UserId       int64     `json:"userId"`
	ImageName    string    `json:"imageName"`
	OrgName      string    `json:"orgName"`
	VersionCount int       `json:"versionCount"`
	UpdateAt     time.Time `json:"updateAt"`
}

type ListContainerImageRsp struct {
	Total  int                     `json:"total"`
	Images []RspContainerImageInfo `json:"images"`
}

// Describe image
type DescribeContainerImageReq struct {
	ImageName string `json:"imageName" validate:"required"`
	OrgName   string `json:"orgName" validate:"required"`
}

type DescribeContainerImageRsp struct {
	Image *imageentity.UserContainerImageInfo `json:"image"`
}

// Upload image
type UploadContainerImageReq struct {
	File    *multipart.FileHeader `form:"file" binding:"required"`
	OrgName string                `form:"orgName" binding:"required"`
}

type UploadContainerImageRsp struct {
}

// Delete image
type DeleteContainerImageReq struct {
	ImageName string `json:"imageName"`
	OrgName   string `json:"orgName"`
}

type DeleteContainerImageRsp struct {
}
