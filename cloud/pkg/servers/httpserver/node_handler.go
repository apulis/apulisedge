// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/node/entity"
	nodeservice "github.com/apulis/ApulisEdge/cloud/pkg/node/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func NodeHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/node")

	group.POST("/createNode", wrapper(CreateEdgeNode))
	group.POST("/listNodes", wrapper(ListEdgeNodes))
	group.POST("/desNode", wrapper(DescribeEgeNode))

}

// Create edge node
type CreateEdgeNodeReq struct {
	proto.ApulisHeader
	NodeName string `json:"nodeName"`
}

type CreateEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// List edge nodes
type ListEdgeNodesReq struct {
	proto.ApulisHeader
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type ListEdgeNodesRsp struct {
	Total int                         `json:"total"`
	Nodes *[]nodeentity.NodeBasicInfo `json:"nodes"`
}

// Describe edge node proto
type DescribeEdgeNodesReq struct {
	proto.ApulisHeader
	NodeName string `json:"nodeName"`
}

type DescribeEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// create edge node
func CreateEdgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent CreateEdgeNodeReq
	var node *nodeentity.NodeBasicInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// create node
	node, err = nodeservice.CreateEdgeNode(reqContent.UserId, reqContent.NodeName)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := CreateEdgeNodesRsp{
		Node: node,
	}
	return SuccessResp(c, &req, data)
}

// list edge nodes
func ListEdgeNodes(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent ListEdgeNodesReq
	var nodes *[]nodeentity.NodeBasicInfo
	var total int

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// list node
	nodes, total, err = nodeservice.ListEdgeNodes(reqContent.UserId, reqContent.PageNum, reqContent.PageSize)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := ListEdgeNodesRsp{
		Total: total,
		Nodes: nodes,
	}
	return SuccessResp(c, &req, data)
}

// describe edge nodes
func DescribeEgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent DescribeEdgeNodesReq
	var nodeInfo *nodeentity.NodeBasicInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// describe node
	nodeInfo, err = nodeservice.DescribeEdgeNode(reqContent.UserId, reqContent.NodeName)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := DescribeEdgeNodesRsp{
		Node: nodeInfo,
	}
	return SuccessResp(c, &req, data)
}
