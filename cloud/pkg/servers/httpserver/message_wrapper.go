// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
	"net/http"
)

var logger = loggers.LogInstance()

type HandlerFunc func(c *gin.Context) error

type APIErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type APISuccessResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

//////////////// Success Message //////////////////
func SuccessResp(c *gin.Context, req *proto.Message, content interface{}) error {
	res := APISuccessResp{
		Code: SUCCESS_CODE,
		Msg:  "OK",
		Data: content,
	}
	rsp := req.NewRespByMessage(req, res)
	c.JSON(http.StatusOK, rsp)
	return nil
}

//////////////// Error Message //////////////////
func (e *APIErrorResp) Error() string {
	return e.Msg
}

func HandleNotFound(c *gin.Context) {
	var req proto.Message

	// no need to handle err here
	_ = c.ShouldBindJSON(&req)
	rsp, _ := ErrorResp(c, &req, NOT_FOUND_ERROR_CODE, http.StatusText(http.StatusNotFound))
	c.JSON(http.StatusNotFound, rsp)

}

func ErrorResp(c *gin.Context, req *proto.Message, code int, msg string) (*proto.Message, *APIErrorResp) {
	errMsg := &APIErrorResp{
		Code: code,
		Msg:  msg,
	}
	rsp := req.NewRespByMessage(req, errMsg)
	return rsp, errMsg
}

func UnAuthorizedError(c *gin.Context, req *proto.Message, msg string) *APIErrorResp {
	rsp, err := ErrorResp(c, req, AUTH_ERROR_CODE, msg)
	c.JSON(http.StatusUnauthorized, rsp)
	return err
}

func ServerError(c *gin.Context, req *proto.Message) *APIErrorResp {
	rsp, err := ErrorResp(c, req, SERVER_ERROR_CODE, http.StatusText(http.StatusInternalServerError))
	c.JSON(http.StatusInternalServerError, rsp)
	return err
}

func NotFoundError(c *gin.Context, req *proto.Message) *APIErrorResp {
	rsp, err := ErrorResp(c, req, NOT_FOUND_ERROR_CODE, http.StatusText(http.StatusNotFound))
	c.JSON(http.StatusNotFound, rsp)
	return err
}

func UnknownError(c *gin.Context, req *proto.Message, msg string) *APIErrorResp {
	rsp, err := ErrorResp(c, req, UNKNOWN_ERROR_CODE, msg)
	c.JSON(http.StatusForbidden, rsp)
	return err
}

func ParameterError(c *gin.Context, req *proto.Message, msg string) *APIErrorResp {
	rsp, err := ErrorResp(c, req, PARAMETER_ERROR_CODE, msg)
	c.JSON(http.StatusBadRequest, rsp)
	return err
}

func AppError(c *gin.Context, req *proto.Message, errCode int, msg string) *APIErrorResp {
	rsp, err := ErrorResp(c, req, errCode, msg)
	c.JSON(http.StatusBadRequest, rsp)
	return err
}

func ServeError(c *gin.Context, req *proto.Message, errCode int, msg string) *APIErrorResp {
	rsp, err := ErrorResp(c, req, errCode, msg)
	c.JSON(http.StatusInternalServerError, rsp)
	return err
}
