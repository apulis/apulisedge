// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	nodeservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func NodeHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/node")
	group.Use(Auth())

	group.POST("/createNode", wrapper(CreateEdgeNode))
	group.POST("/listNode", wrapper(ListEdgeNodes))
	group.POST("/desNode", wrapper(DescribeEgeNode))
	group.POST("/deleteNode", wrapper(DeleteEdgeNode))
}

// create edge node
func CreateEdgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.CreateEdgeNodeReq
	var node *nodeentity.NodeBasicInfo

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// create node
	node, err = nodeservice.CreateEdgeNode(reqContent.UserId, reqContent.UserName, reqContent.NodeName)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := nodemodule.CreateEdgeNodeRsp{
		Node: node,
	}
	return SuccessResp(c, &req, data)
}

// list edge nodes
func ListEdgeNodes(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.ListEdgeNodesReq
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

	data := nodemodule.ListEdgeNodesRsp{
		Total: total,
		Nodes: nodes,
	}
	return SuccessResp(c, &req, data)
}

// describe edge nodes
func DescribeEgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.DescribeEdgeNodesReq
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

	data := nodemodule.DescribeEdgeNodesRsp{
		Node: nodeInfo,
	}
	return SuccessResp(c, &req, data)
}

// delete edge node
func DeleteEdgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.DeleteEdgeNodeReq

	if err = c.ShouldBindJSON(&req); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	if err := mapstructure.Decode(req.Content.(map[string]interface{}), &reqContent); err != nil {
		return ParameterError(c, &req, err.Error())
	}

	// TODO validate reqContent

	// delete application
	err = nodeservice.DeleteEdgeNode(&reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}
