package httpserver

import (
	"net/http"

	"github.com/apulis/ApulisEdge/cloud/pkg/domain/authentication"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

// AuthenticationHandlerRoutes join authentication module in server
func AuthenticationHandlerRoutes(r *gin.Engine) {

	group := r.Group("/apulisEdge/api/authentication")
	group.Use(Auth())

	group.GET("/test", wrapper(authenticationTest))
}

func authenticationTest(c *gin.Context) error {
	var req proto.Message
	data := "success"
	return SuccessResp(c, &req, data)
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req proto.Message
		authenticated := authentication.Auth(c)
		if authenticated.Result {

		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, UnAuthorizedError(c, &req, "authentication fail"))
		}
		c.Next()

	}
}
