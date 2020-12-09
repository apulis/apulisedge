package authentication

import (
	"fmt"

	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authMethod func(c *gin.Context) AuthResult
var logger = loggers.LogInstance()

// Auth offer multi authentication features
func Auth(c *gin.Context) AuthResult {
	authType := viper.GetStringMap("authentication")["type"]
	if authType == "none" {
		return NoAuth
	}
	if authType == "JWT" {
		authMethod = jwtAuthtication
		jwtSecretKey = viper.GetStringMap("authentication")["key"].(string)
	} else if authType == "basic" {
		username = viper.GetViper().GetStringMap("authentication")["username"].(string)
		password = viper.GetViper().GetStringMap("authentication")["password"].(string)
		authMethod = basicAuthentication
	} else {
		panic(fmt.Errorf("unsupport authentication metho: %s", authType))
	}
	return authMethod(c)
}

func Sign() {
}
