// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// List deploy
type ListEdgeAppDeployReq struct {
	AppName  string `json:"appName"`
	Version  string `json:"version"`
	PageNum  int    `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListEdgeAppDeployRsp struct {
	Total      int
	AppDeploys *[]appentity.ApplicationDeployInfo `json:"appDeploys"`
}

// List node deploy
type ListNodeDeployReq struct {
	Name     string `json:"name" validate:"required"`
	PageNum  int    `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListNodeDeployRsp struct {
	Total      int
	AppDeploys *[]appentity.ApplicationDeployInfo `json:"appDeploys"`
}

// Deploy edge application
type DeployEdgeApplicationReq struct {
	AppName   string   `json:"appName" validate:"required"`
	NodeNames []string `json:"nodeNames" validate:"required"`
	Version   string   `json:"version" validate:"required"`
}

type DeployEdgeApplicationRsp struct {
}

// undeploy edge application
type UnDeployEdgeApplicationReq struct {
	AppName   string   `json:"appName" validate:"required"`
	NodeNames []string `json:"nodeNames" validate:"required"`
	Version   string   `json:"version" validate:"required"`
}

type UnDeployEdgeApplicationRsp struct {
}
