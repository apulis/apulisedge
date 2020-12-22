// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	nodeservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/service"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

func NodeHandlerRoutes(r *gin.Engine) {
	group := r.Group("/apulisEdge/api/node")

	// add authentication
	group.Use(Auth())

	group.POST("/createNode", wrapper(CreateEdgeNode))
	group.POST("/listNode", wrapper(ListEdgeNodes))
	group.POST("/desNode", wrapper(DescribeEdgeNode))
	group.POST("/deleteNode", wrapper(DeleteEdgeNode))
	group.POST("/scripts", wrapper(GetInstallScripts))
	group.POST("/listType", wrapper(ListEdgeNodeType))
}

// @Summary create edge node
// @Description create edge node
// @Tags node
// @Accept json
// @Produce json
// @Param param body node.CreateEdgeNodeReq true "userId:user id, userName: user name"
// @Success 200 {object} APISuccessResp{data=node.CreateEdgeNodeRsp}
// @Failure 400 {object} APIErrorResp
// @Router /createNode [post]
func CreateEdgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.CreateEdgeNodeReq
	var node *nodeentity.NodeBasicInfo

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// create node
	node, err = nodeservice.CreateEdgeNode(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := nodemodule.CreateEdgeNodeRsp{
		Node: node,
	}
	return SuccessResp(c, &req, data)
}

// @Summary list edge nodes
// @Description list edge nodes
// @Tags node
// @Accept json
// @Produce json
// @Param param body node.ListEdgeNodesReq true "userId:user id, userName: user name"
// @Success 200 {object} APISuccessResp{data=node.ListEdgeNodesRsp} "code:0, msg:OK"
// @Failure 400 {object} APIErrorResp "code:30000, msg:db error"
// @Router /listNode [post]
func ListEdgeNodes(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.ListEdgeNodesReq
	var nodes *[]nodeentity.NodeBasicInfo
	var total int

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list node
	nodes, total, err = nodeservice.ListEdgeNodes(*userInfo, &reqContent)
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
func DescribeEdgeNode(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.DescribeEdgeNodesReq
	var nodeInfo *nodeentity.NodeBasicInfo

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// describe node
	nodeInfo, err = nodeservice.DescribeEdgeNode(*userInfo, &reqContent)
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

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// delete application
	err = nodeservice.DeleteEdgeNode(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, "OK")
}

func GetInstallScripts(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.GetInstallScriptReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// get cluster
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		logger.Infof("GetInstallScripts, can`t find cluster %d", userInfo.ClusterId)
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	script, err := nodeservice.GetInstallScripts(&reqContent, clu.Domain, clu.HarborAddress, clu.DownloadAddress)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := nodemodule.GetInstallScriptRsp{
		Script: script,
	}
	return SuccessResp(c, &req, data)
}

func ListEdgeNodeType(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.ListEdgeNodeTypeReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list type
	tys, err := nodeservice.ListEdgeNodeType(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := nodemodule.ListEdgeNodeTypeRsq{
		Types: tys,
	}
	return SuccessResp(c, &req, data)
}
