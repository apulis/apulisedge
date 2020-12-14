package httpserver

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"net/http"

	"github.com/apulis/ApulisEdge/cloud/pkg/domain/authentication"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/gin-gonic/gin"
)

var authenticator authentication.Authenticator

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

func InitAuth(config *configs.EdgeCloudConfig) error {
	var err error
	authenticator, err = authentication.GetAuthenticator(config)
	if err != nil {
		return err
	}

	logger.Infof("Init auth succ, now use authType = %s", config.Authentication.AuthType)
	return nil
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		r := c.Request
		auth := r.Header.Get("Authorization")
		logger.Debugf("Auth header = %s", auth)
		if len(auth) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, NoBodyUnAuthorizedError("Cannot authorize"))
			c.Next()
			return
		}

		data, err := authenticator.AuthMethod(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, NoBodyUnAuthorizedError(err.Error()))
		} else {
			c.Set("clusterId", data.ClusterId)
			c.Set("groupId", data.GroupId)
			c.Set("userId", data.UserId)
		}

		c.Next()
	}
}
