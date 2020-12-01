// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	nodeentity "github.com/apulis/ApulisEdge/cloud/node/entity"
	nodeservice "github.com/apulis/ApulisEdge/cloud/node/service"
	proto "github.com/apulis/ApulisEdge/cloud/protocol"
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
	NodeName string `json:"nodeName"`
}

type CreateEdgeNodesRsp struct {
	Node *nodeentity.NodeBasicInfo `json:"node"`
}

// List edge nodes
type ListEdgeNodesReq struct {
}

type ListEdgeNodesRsp struct {
	Nodes []*nodeentity.NodeBasicInfo `json:"nodes"`
	//Total     int                         `json:"total"`
	//TotalPage int                         `json:"totalPage"`
}

// Describe edge node proto
type DescribeEdgeNodesReq struct {
	NodeName string `json:"nodeName"`
}

type DescribeEdgeNodesRsp struct {
	Node *nodeentity.NodeDetailInfo `json:"node"`
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

	// TODO validate NodeName

	// create node
	node, err = nodeservice.CreateEdgeNode(reqContent.NodeName)
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
	var nodes []*nodeentity.NodeBasicInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// list node
	nodes, err = nodeservice.ListEdgeNodes()
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := ListEdgeNodesRsp{
		Nodes: nodes,
	}
	return SuccessResp(c, &req, data)
}

// describe edge nodes
func DescribeEgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent DescribeEdgeNodesReq
	var nodeInfo *nodeentity.NodeDetailInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate NodeName

	// describe node
	nodeInfo, err = nodeservice.DescribeEdgeNode(reqContent.NodeName)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := DescribeEdgeNodesRsp{
		Node: nodeInfo,
	}
	return SuccessResp(c, &req, data)
}
