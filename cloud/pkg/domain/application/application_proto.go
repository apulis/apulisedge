// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// List edge application
type ListEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	PageNum            int `json:"pageNum"`
	PageSize           int `json:"pageSize"`
}

type ListEdgeApplicationRsp struct {
	Total int
	Apps  *[]appentity.ApplicationBasicInfo `json:"apps"`
}

// Create edge application
type CreateEdgeApplicationReq struct {
	proto.ApulisHeader    `mapstructure:",squash"`
	AppName               string  `json:"appName"`
	ArchType              int     `json:"archType"`
	Version               string  `json:"version"`
	ContainerImage        string  `json:"containerImage"`
	ContainerImageVersion string  `json:"containerImageVersion"`
	ContainerImagePath    string  `json:"containerImagePath"`
	CpuQuota              float32 `json:"cpuQuota"`
	MemoryQuota           int     `json:"memoryQuota"`
}

type CreateEdgeApplicationRsp struct {
	Application *appentity.ApplicationBasicInfo `json:"application"`
}

// Delete edge application
type DeleteEdgeApplicationReq struct {
	proto.ApulisHeader `mapstructure:",squash"`
	AppName            string `json:"appName"`
	Version            string `json:"version"`
}

type DeleteEdgeApplicationRsp struct {
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
	NamespaceName      string `json:"namespaceName"`
	Version            string `json:"version"`
	RestartPolicy      int    `json:"restartPolicy"`
	ContainerPort      int    `json:"containerPort"`
	PortMapping        struct {
		Enable   bool `json:"enable"`
		HostPort int  `json:"hostPort"`
	} `json:"portMapping"`
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
