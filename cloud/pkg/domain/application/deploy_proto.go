// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// List deploy
type ListEdgeAppDeployReq struct {
	AppName  string `json:"appName" validate:"required"`
	Version  string `json:"version" validate:"required"`
	PageNum  int    `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListEdgeAppDeployRsp struct {
	Total      int
	AppDeploys *[]appentity.ApplicationDeployInfo `json:"appDeploys"`
}

// List node deploy
type ListNodeDeployReq struct {
	Name     string `json:"name" validate:"required"`
	PageNum  int    `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListNodeDeployRsp struct {
	Total      int
	AppDeploys *[]appentity.ApplicationDeployInfo `json:"appDeploys"`
}

// List node can deploy
type ListNodeCanDeployReq struct {
	AppName  string `json:"appName" validate:"required"`
	Version  string `json:"version" validate:"required"`
	PageNum  int    `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListNodeCanDeployRsp struct {
	Total int
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// List node can update
type ListNodeCanUpdateReq struct {
	AppName  string `json:"appName" validate:"required"`
	Version  string `json:"version" validate:"required"`
	PageNum  int    `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListNodeCanUpdateRsp struct {
	Total int
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// Deploy edge application
type DeployEdgeApplicationReq struct {
	AppName   string   `json:"appName" validate:"required"`
	NodeNames []string `json:"nodeNames" validate:"required"`
	Version   string   `json:"version" validate:"required"`
}

type DeployEdgeApplicationRsp struct {
}

// Update deploy edge application
type UpdateDeployEdgeApplicationReq struct {
	AppName       string   `json:"appName" validate:"required"`
	NodeNames     []string `json:"nodeNames" validate:"required"`
	TargetVersion string   `json:"targetVersion" validate:"required"`
}

type UpdateDeployEdgeApplicationRsp struct {
}

// undeploy edge application
type UnDeployEdgeApplicationReq struct {
	AppName   string   `json:"appName" validate:"required"`
	NodeNames []string `json:"nodeNames" validate:"required"`
	Version   string   `json:"version" validate:"required"`
}

type UnDeployEdgeApplicationRsp struct {
}
