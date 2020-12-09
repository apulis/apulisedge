package httpserver

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/authentication"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

// AuthenticationHandlerRoutes join authentication module in server
func AuthenticationHandlerRoutes(r *gin.Engine) {

	group := r.Group("/apulisEdge/api/authentication")
	group.Use(authentication.Auth())

	group.GET("/test", wrapper(authenticationTest))
}

func authenticationTest(c *gin.Context) error {
	var req proto.Message
	data := "success"
	return SuccessResp(c, &req, data)
}
