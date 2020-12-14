// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// Create edge node
type CreateEdgeNodeReq struct {
	protocol.ApulisHeader `mapstructure:",squash"`
	Name                  string `json:"name"`
}

type CreateEdgeNodeRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// List edge nodes
type ListEdgeNodesReq struct {
	protocol.ApulisHeader `mapstructure:",squash"`
	PageNum               int `json:"pageNum"`
	PageSize              int `json:"pageSize"`
}

type ListEdgeNodesRsp struct {
	Total int                         `json:"total"`
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// Describe edge node protocol
type DescribeEdgeNodesReq struct {
	protocol.ApulisHeader `mapstructure:",squash"`
	Name                  string `json:"name"`
}

type DescribeEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// Delete edge node
type DeleteEdgeNodeReq struct {
	protocol.ApulisHeader `mapstructure:",squash"`
	Name                  string `json:"name"`
}

type DeleteEdgeNodeRsp struct {
}
