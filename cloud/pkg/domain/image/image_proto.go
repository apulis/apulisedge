// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"time"
)

// List image
type ListContainerImageReq struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
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

// Upload image
type UploadContainerImageReq struct {
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

// List image version
type ListContainerImageVersionReq struct {
	ImageName string `json:"imageName"`
	OrgName   string `json:"orgName"`
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
}

type ListContainerImageVersionRsp struct {
	Total  int                                          `json:"total"`
	Images *[]imageentity.UserContainerImageVersionInfo `json:"imageVersions"`
}

// Delete image version
type DeleteContainerImageVersionReq struct {
	ImageName    string `json:"imageName"`
	OrgName      string `json:"orgName"`
	ImageVersion string `json:"imageVersion"`
}

type DeleteContainerImageVersionRsp struct {
}

// List image org
type ListContainerImageOrgReq struct {
	OrgName  string `json:"orgName"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

type ListContainerImageOrgRsp struct {
	Total  int                              `json:"total"`
	Images *[]imageentity.ContainerImageOrg `json:"imageOrgs"`
}

// Delete image org
type DeleteContainerImageOrgReq struct {
	OrgName string `json:"orgName"`
}

type DeleteContainerImageOrgRsp struct {
}
