// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// Create edge application ////////////////////////////////////////////////////////
type CreateEdgeApplicationReq struct {
	AppName               string           `json:"appName"`
	FunctionType          string           `json:"functionType"`
	Description           string           `json:"description"`
	ArchType              []string         `json:"archType" validate:"required"`
	Version               string           `json:"version"`
	OrgName               string           `json:"orgName"`
	ContainerImage        string           `json:"containerImage"`
	ContainerImageVersion string           `json:"containerImageVersion"`
	CpuQuota              float32          `json:"cpuQuota"`
	MaxCpuQuota           float32          `json:"maxCpuQuota"`
	MemoryQuota           float32          `json:"memoryQuota"`
	MaxMemoryQuota        float32          `json:"maxMemoryQuota"`
	RestartPolicy         string           `json:"restartPolicy" validate:"required,oneof=Always OnFailure Never"`
	Network               CreateAppNetwork `json:"network" validate:"required"`
}

type CreateAppNetwork struct {
	Type         string        `json:"type" validate:"required,oneof=Host PortMapping"`
	PortMappings []PortMapping `json:"portMappings"`
}

type PortMapping struct {
	ContainerPort int `json:"containerPort" validate:"required"`
	HostPort      int `json:"hostPort" validate:"required"`
}

type CreateEdgeApplicationRsp struct {
	AppCreated     string `json:"appCreated"`
	VersionCreated string `json:"versionCreated"`
}

// List edge application ////////////////////////////////////////////////////////
type ListEdgeApplicationReq struct {
	AppType  string `json:"appType" validate:"oneof=UserDefine System All"`
	PageNum  int    `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int    `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListEdgeApplicationRsp struct {
	Total int                               `json:"total"`
	Apps  *[]appentity.ApplicationBasicInfo `json:"apps"`
}

// Describe edge application ///////////////////////////////////////////////////////
type DescribeEdgeApplicationReq struct {
	AppName string `json:"name"`
}

type DescribeEdgeApplicationRsp struct {
	App *appentity.ApplicationBasicInfo `json:"app"`
}

// Delete edge application ///////////////////////////////////////////////////////
type DeleteEdgeApplicationReq struct {
	AppName string `json:"appName"`
}

type DeleteEdgeApplicationRsp struct {
}
