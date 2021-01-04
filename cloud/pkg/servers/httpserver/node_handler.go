// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"path"
	"time"

	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
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
	group.POST("/listArchType", wrapper(ListArchType))

	group.POST("/batchcsv", wrapper(UploadNodeFile))
	group.POST("/batchnode", wrapper(AddBatchNode))
	group.POST("/batch", wrapper(ConfirmNodesBatch))
	group.GET("/batch", wrapper(GetTempNodesBatchList))
	group.DELETE("/batch", wrapper(DeleteNodeOfBatch))
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

	script, err := nodeservice.GetInstallScripts(*userInfo, &reqContent)
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

	data := nodemodule.ListEdgeNodeTypeRsp{
		Types: tys,
	}
	return SuccessResp(c, &req, data)
}

func ListArchType(c *gin.Context) error {
	var err error
	var req proto.Message
	var reqContent nodemodule.ListArchTypeReq

	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}

	// list type
	tys, err := nodeservice.ListArchType(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	data := nodemodule.ListArchTypeRsp{
		Types: tys,
	}
	return SuccessResp(c, &req, data)
}

func UploadNodeFile(c *gin.Context) error {
	var req proto.Message
	data := "success"
	file, err := c.FormFile("data")
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}
	csvPath := path.Join(nodemodule.CSVSavePath, time.Now().Format(time.RFC3339)+".csv")
	err = c.SaveUploadedFile(file, csvPath)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}
	// get user info, user info comes from authentication
	userInfo := proto.ApulisHeader{}
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}
	err = nodeservice.AnalyzeCSV(userInfo, csvPath)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, data)
}

func AddBatchNode(c *gin.Context) error {
	var req proto.Message
	var reqContent nodemodule.CreateNodeOfBatchReq
	data := "success"
	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}
	err := nodeservice.AddOneRecord(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, data)
}

func GetTempNodesBatchList(c *gin.Context) error {
	var req proto.Message
	var reqContent nodemodule.ListNodeOfBatchReq
	// get user info, user info comes from authentication
	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return AppError(c, &req, APP_ERROR_CODE, errRsp.Error())
	}
	batchList, err := nodeservice.ListBatchList(userInfo, reqContent.PageSize, reqContent.PageNum)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}
	rspContent := nodemodule.ListNodeOfBatchRsp{
		Status:   "loaded",
		NodeList: *batchList,
	}

	return SuccessResp(c, &req, rspContent)
}

func ConfirmNodesBatch(c *gin.Context) error {
	var req proto.Message
	data := "success"
	userInfo := proto.ApulisHeader{}
	var err error
	userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, err = GetUserInfo(c)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}
	go nodeservice.UpdateBatch(userInfo)

	return SuccessResp(c, &req, data)
}

func DeleteNodeOfBatch(c *gin.Context) error {
	var req proto.Message
	var reqContent nodemodule.DeleteNodeOfBatchReq
	data := "success"
	userInfo, errRsp := PreHandler(c, &req, &reqContent)
	if errRsp != nil {
		return errRsp
	}
	err := nodeservice.DeleteNodeOfBatch(*userInfo, &reqContent)
	if err != nil {
		return AppError(c, &req, APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, &req, data)
}
