// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// Create node group
type CreateNodeGroupReq struct {
	GroupName string `json:"name" validate:"required"`
}

type CreateNodeGroupRsp struct {
	Group *nodeentity.NodeGroupInfo `json:"node"`
}

// List group
type ListNodeGroupReq struct {
	PageNum  int `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListNodeGroupRsp struct {
	Total  int                         `json:"total"`
	Groups *[]nodeentity.NodeGroupInfo `json:"groups"`
}

// Describe group
type DescribeNodeGroupReq struct {
	GroupName string `json:"groupName" validate:"required"`
}

type DescribeNodeGroupRsp struct {
	Group *nodeentity.NodeGroupInfo `json:"group"`
}

// Delete node group
type DeleteNodeGroupReq struct {
	GroupName string `json:"groupName" validate:"required"`
}

type DeleteNodeGroupRsp struct {
}
