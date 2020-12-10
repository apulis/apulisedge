package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	AuthMethod(*gin.Context) AuthResult
	initCertificate()
}
type AuthResult struct {
	Result    bool
	AuthError error
}

var NoAuth = AuthResult{
	Result:    true,
	AuthError: nil,
}

var NoAuthHeadError = AuthResult{
	Result:    false,
	AuthError: errors.New("can't found authentication in header"),
}

var BasicAuthFailError = AuthResult{
	Result:    false,
	AuthError: errors.New("basic auth fails"),
}

var BasicAuthDecodeFailError = AuthResult{
	Result:    false,
	AuthError: errors.New("basic auth decode fails"),
}

var JWTAuthFailError = AuthResult{
	Result:    false,
	AuthError: errors.New("JWT auth fail"),
}

func newAuthResult(result bool, authError error) AuthResult {
	return AuthResult{
		Result:    result,
		AuthError: authError,
	}
}
