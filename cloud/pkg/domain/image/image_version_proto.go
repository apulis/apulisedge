// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// List image version
type ListContainerImageVersionReq struct {
	ImageName string `json:"imageName"`
	OrgName   string `json:"orgName"`
	PageNum   int    `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize  int    `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListContainerImageVersionRsp struct {
	Total         int                                          `json:"total"`
	ImageVersions *[]imageentity.UserContainerImageVersionInfo `json:"imageVersions"`
}

// Describe image version
type DescribeContainerImageVersionReq struct {
	ImageName    string `json:"imageName" validate:"required"`
	OrgName      string `json:"orgName" validate:"required"`
	ImageVersion string `json:"imageVersion" validate:"required"`
}

type DescribeContainerImageVersionRsp struct {
	ImageVersion *imageentity.UserContainerImageVersionInfo `json:"imageVersion"`
}

// Delete image version
type DeleteContainerImageVersionReq struct {
	ImageName    string `json:"imageName" validate:"required"`
	OrgName      string `json:"orgName" validate:"required"`
	ImageVersion string `json:"imageVersion" validate:"required"`
}

type DeleteContainerImageVersionRsp struct {
}
