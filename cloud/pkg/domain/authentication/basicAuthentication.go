package authentication

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var username string
var password string

type basicAuthentication struct {
}

func (basicAuthenticator basicAuthentication) AuthMethod(c *gin.Context) AuthResult {
	authenticated := AuthResult{
		Result:    false,
		AuthError: nil,
	}

	r := c.Request
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return NoAuthHeadError
	}
	tokenString := strings.Fields(auth)[1]

	basicAuthString, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return newAuthResult(false, err)
	}

	reqUsername := strings.Split(string(basicAuthString), ":")[0]
	reqPassword := strings.Split(string(basicAuthString), ":")[1]

	if reqUsername == username && reqPassword == password {
		authenticated.Result = true
	} else {
		return newAuthResult(false, errors.New("user name and password not valid"))
	}

	return authenticated
}

func (basicAuthenticator basicAuthentication) initCertificate() {
	username = viper.GetViper().GetStringMap("authentication")["username"].(string)
	password = viper.GetViper().GetStringMap("authentication")["password"].(string)
}
