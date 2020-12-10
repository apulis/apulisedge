// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// Create edge application
type CreateEdgeApplicationReq struct {
	proto.ApulisHeader    `mapstructure:",squash"`
	AppName               string  `json:"appName"`
	FunctionType          string  `json:"functionType"`
	Description           string  `json:"description"`
	ArchType              string  `json:"archType"`
	Version               string  `json:"version"`
	ContainerImage        string  `json:"containerImage"`
	ContainerImageVersion string  `json:"containerImageVersion"`
	ContainerImagePath    string  `json:"containerImagePath"`
	CpuQuota              float32 `json:"cpuQuota"`
	MaxCpuQuota           float32 `json:"maxCpuQuota"`
	MemoryQuota           float32 `json:"memoryQuota"`
	MaxMemoryQuota        float32 `json:"maxMemoryQuota"`
}

type CreateEdgeApplicationRsp struct {
	AppCreated     string `json:"appCreated"`
	VersionCreated string `json:"versionCreated"`
}

// List edge application
type ListEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	PageNum            int `json:"pageNum"`
	PageSize           int `json:"pageSize"`
}

type ListEdgeApplicationRsp struct {
	Total int                               `json:"total"`
	Apps  *[]appentity.ApplicationBasicInfo `json:"apps"`
}

// List edge application version
type ListEdgeApplicationVersionReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	PageNum            int    `json:"pageNum"`
	PageSize           int    `json:"pageSize"`
}

type ListEdgeApplicationVersionRsp struct {
	Total       int                                 `json:"total"`
	AppVersions *[]appentity.ApplicationVersionInfo `json:"appVersions"`
}

// Delete edge application
type DeleteEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
}

type DeleteEdgeApplicationRsp struct {
}

// Delete edge application version
type DeleteEdgeApplicationVersionReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	Version            string `json:"version"`
}

type DeleteEdgeApplicationVersionRsp struct {
}

// List edge application
type ListEdgeAppDeployReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	Version            string `json:"version"`
	PageNum            int    `json:"pageNum"`
	PageSize           int    `json:"pageSize"`
}

type ListEdgeAppDeployRsp struct {
	Total      int
	AppDeploys *[]appentity.ApplicationDeployInfo `json:"appDeploys"`
}

// Deploy edge application
type DeployEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	NodeName           string `json:"nodeName"`
	Version            string `json:"version"`
}

type DeployEdgeApplicationRsp struct {
}

// undeploy edge application
type UnDeployEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	NodeName           string `json:"nodeName"`
	Version            string `json:"version"`
}

type UnDeployEdgeApplicationRsp struct {
}
