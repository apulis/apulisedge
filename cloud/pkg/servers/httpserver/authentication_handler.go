package httpserver

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"net/http"
	"strings"

	"github.com/apulis/ApulisEdge/cloud/pkg/domain/authentication"
	"github.com/gin-gonic/gin"
)

var authenticator authentication.Authenticator

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
		authHeader := r.Header.Get("Authorization")
		logger.Debugf("Auth header = %s", authHeader)
		if len(authHeader) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, NoBodyUnAuthorizedError("Cannot authorize"))
			c.Next()
			return
		}

		auth := strings.Fields(authHeader)[1]
		data, err := authenticator.AuthMethod(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, NoBodyUnAuthorizedError(err.Error()))
			c.Next()
			return
		}

		// TODO check clusterId/groupId/userId
		if data.ClusterId < 0 || data.GroupId < 0 || data.UserId < 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, NoBodyUnAuthorizedError(ErrInvalidUserInfo.Error()))
			c.Next()
			return
		}

		c.Set("clusterId", data.ClusterId)
		c.Set("groupId", data.GroupId)
		c.Set("userId", data.UserId)
		c.Next()
	}
}
