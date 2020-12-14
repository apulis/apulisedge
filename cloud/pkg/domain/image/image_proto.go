// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"time"
)

// List edge nodes
type ListContainerImageReq struct {
	protocol.ApulisHeader `mapstructure:",squash"`
	PageNum               int `json:"pageNum"`
	PageSize              int `json:"pageSize"`
}

type RspImageInfo struct {
	ClusterId    int64     `json:"clusterId"`
	GroupId      int64     `json:"groupId"`
	UserId       int64     `json:"userId"`
	ImageName    string    `json:"imageName"`
	OrgName      string    `json:"orgName"`
	VersionCount int       `json:"versionCount"`
	UpdateAt     time.Time `json:"updateAt"`
}

type ListContainerImageRsp struct {
	Total  int            `json:"total"`
	Images []RspImageInfo `json:"images"`
}
