package authentication

import (
	"fmt"

	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authenticator Authenticator
var logger = loggers.LogInstance()

// Auth offer multi authentication features
func Auth(c *gin.Context) AuthResult {
	authType := viper.GetStringMap("authentication")["type"]
	if authType == "none" || authType == nil {
		return NoAuth
	}
	if authType == "JWT" {
		authenticator = new(jwtAuthtication)
	} else if authType == "basic" {
		authenticator = new(basicAuthentication)
	} else {
		panic(fmt.Errorf("unsupport authentication metho: %s", authType))
	}

	authenticator.initCertificate()
	return authenticator.AuthMethod(c)
}

func Sign() {
}
