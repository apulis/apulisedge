// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeentity

type NodeBasicInfo struct {
	Name             string `json:"name"   binding:"required"`
	Status           string `json:"status" binding:"required"`
	Roles            string `json:"roles" binding:"required"`
	ContainerRuntime string `json:"runtime" binding:"required"`
	OsImage          string `json:"osImage" binding:"required"`
	ProviderId       string `json:"providerId"`
	InterIp          string `json:"interIp"`
	OuterIp          string `json:"outerIp"`
	CreateTime       string `json:"createTime"`
}

type NodeDetailInfo struct {
	Name             string `json:"name"   binding:"required"`
	Status           string `json:"status" binding:"required"`
	Roles            string `json:"roles" binding:"required"`
	ContainerRuntime string `json:"runtime" binding:"required"`
	OsImage          string `json:"osImage" binding:"required"`
	ProviderId       string `json:"providerId"`
	InterIp          string `json:"interIp"`
	OuterIp          string `json:"outerIp"`
	CreateTime       string `json:"createTime"`
}
