package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authMethod gin.HandlerFunc

// Auth offer multi authentication features
func Auth() gin.HandlerFunc {
	authType := viper.GetStringMap("authentication")["type"]
	if authType == "JWT" {
		authMethod = jwtAuthtication
	} else {
		panic(fmt.Errorf("unsupport authentication metho: %s", authType))
	}
	return authMethod
}

func Sign() {
}
