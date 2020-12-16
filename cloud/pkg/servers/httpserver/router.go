// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	_ "github.com/apulis/ApulisEdge/cloud/pkg/docs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title ApulisEdge Cloud API
// @version alpha
// @description ApulisEdge cloud server.
func NewRouter() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(cors.Default())

	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	r.Use(loggers.GinLogger(logger))
	r.Use(gin.Recovery())

	NodeHandlerRoutes(r)
	ApplicationHandlerRoutes(r)
	ImageHandlerRoutes(r)
	return r
}
