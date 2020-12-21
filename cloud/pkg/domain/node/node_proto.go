// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
)

// TODO add param validate, like node.ListEdgeNodesReq

// Create edge node
type CreateEdgeNodeReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateEdgeNodeRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// List edge nodes
type ListEdgeNodesReq struct {
	PageNum  int `json:"pageNum" validate:"gte=1,lte=1000"`
	PageSize int `json:"pageSize" validate:"gte=1,lte=1000"`
}

type ListEdgeNodesRsp struct {
	Total int                         `json:"total"`
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// Describe edge node protocol
type DescribeEdgeNodesReq struct {
	Name string `json:"name"`
}

type DescribeEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// Delete edge node
type DeleteEdgeNodeReq struct {
	Name string `json:"name"`
}

type DeleteEdgeNodeRsp struct {
}

type GetInstallScriptReq struct {
	Arch string `json:"arch"`
}

type GetInstallScriptRsp struct {
	Script string `json:"script"`
}

/////////// node type ///////////////
// List node type
type ListEdgeNodeTypeReq struct {
}

type ListEdgeNodeTypeRsq struct {
	Types []string `json:"types"`
}
