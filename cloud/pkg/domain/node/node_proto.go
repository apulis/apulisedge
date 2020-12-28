// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// Create edge node
type CreateEdgeNodeReq struct {
	Name     string `json:"name" validate:"required"`
	NodeType string `json:"nodeType" validate:"required"`
}

type CreateEdgeNodeRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// List edge nodes
type ListEdgeNodesReq struct {
	PageNum  int `json:"pageNum" validate:"required,gte=1,lte=1000"`
	PageSize int `json:"pageSize" validate:"required,gte=1,lte=1000"`
}

type ListEdgeNodesRsp struct {
	Total int                         `json:"total"`
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// Describe edge node protocol
type DescribeEdgeNodesReq struct {
	Name string `json:"name" validate:"required"`
}

type DescribeEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// Delete edge node
type DeleteEdgeNodeReq struct {
	Name string `json:"name" validate:"required"`
}

type DeleteEdgeNodeRsp struct {
}

type DeleteNodeOfBatchReq struct {
	Name string `json:"name" validate:"required"`
}
type GetInstallScriptReq struct {
	Name string `json:"name" validate:"required"`
	Arch string `json:"arch" validate:"required"`
}

type GetInstallScriptRsp struct {
	Script string `json:"script"`
}

/////////// node type ///////////////
// List node type
type ListEdgeNodeTypeReq struct {
}

type ListEdgeNodeTypeRsp struct {
	Types []string `json:"types"`
}

/////////// arch type ///////////////
// List arch type
type ListArchTypeReq struct {
}

type ListArchTypeRsp struct {
	Types []string `json:"types"`
}
