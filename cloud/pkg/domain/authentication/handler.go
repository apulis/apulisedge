package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authenticator Authenticator
var authenticatorMap = map[string]Authenticator{
	"JWT":   new(jwtAuthtication),
	"basic": new(basicAuthentication),
}

// Auth offer multi authentication features
// If you want to add a authenticator, follow steps:
// 1. define your authenticator struct, and implements "Authenticator" interface
// 2. add your config key pair in authenticatorMap
func Auth(c *gin.Context) AuthResult {
	authType := viper.GetStringMap("authentication")["type"]
	if authType == "none" || authType == nil {
		return NoAuth
	}

	authenticator := authenticatorMap[authType.(string)]
	if authenticator == nil {
		return NotSupportAuth
	}

	authenticator.initCertificate()
	return authenticator.AuthMethod(c)
}

func IsSupport(authName string) bool {
	_, ok := authenticatorMap[authName]
	return ok
}
