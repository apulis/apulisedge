package authentication

import (
	"fmt"

	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authMethod gin.HandlerFunc
var logger = loggers.LogInstance()

// Auth offer multi authentication features
func Auth() gin.HandlerFunc {
	authType := viper.GetStringMap("authentication")["type"]
	if authType == "JWT" {
		authMethod = jwtAuthtication
		jwtSecretKey = viper.GetStringMap("authentication")["key"].(string)
	} else {
		panic(fmt.Errorf("unsupport authentication metho: %s", authType))
	}
	return authMethod
}

func Sign() {
}
