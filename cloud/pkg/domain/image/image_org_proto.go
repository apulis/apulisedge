// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// Create image org
type CreateContainerImageOrgReq struct {
	OrgName string `json:"orgName"`
}

type CreateContainerImageOrgRsp struct {
	Org *imageentity.ContainerImageOrg `json:"org"`
}

// List image org
type ListContainerImageOrgReq struct {
	PageNum  int `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListContainerImageOrgRsp struct {
	Total     int                              `json:"total"`
	ImageOrgs *[]imageentity.ContainerImageOrg `json:"imageOrgs"`
}

// Delete image org
type DeleteContainerImageOrgReq struct {
	OrgName string `json:"orgName"`
}

type DeleteContainerImageOrgRsp struct {
}
