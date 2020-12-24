// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// List edge application version
type ListEdgeApplicationVersionReq struct {
	AppName  string `json:"appName"`
	PageNum  int    `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListEdgeApplicationVersionRsp struct {
	Total       int                                 `json:"total"`
	AppVersions *[]appentity.ApplicationVersionInfo `json:"appVersions"`
}

/////// Describe edge application version
type DescribeEdgeAppVersionReq struct {
	AppName string `json:"name" validate:"required"`
	Version string `json:"version" validate:"required"`
}

type DescribeEdgeAppVersionRsp struct {
	AppVersion *appentity.ApplicationVersionInfo `json:"appVersion"`
}

///// publish application version
type PublishEdgeApplicationVersionReq struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

type PublishEdgeApplicationVersionRsp struct {
}

///// offline application version
type OfflineEdgeApplicationVersionReq struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

type OfflineEdgeApplicationVersionRsp struct {
}

// Delete edge application version
type DeleteEdgeApplicationVersionReq struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

type DeleteEdgeApplicationVersionRsp struct {
}
