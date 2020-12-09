package authentication

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
)

var username string
var password string

func basicAuthentication(c *gin.Context) AuthResult {
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
	}

	return authenticated
}
