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
