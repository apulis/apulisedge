package authentication

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecretKey string

func jwtAuthtication(c *gin.Context) {

	r := c.Request
	auth := r.Header.Get("Authorization")
	tokenString := strings.Fields(auth)[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if token.Valid {
		logger.Infoln("token valid")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			logger.Errorln("token format error")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			logger.Errorln("token expired or nor valid yet")
		} else {
			logger.Errorln("Couldn't handle this token:", err)
		}
	} else {
		logger.Errorln("Couldn't handle this token:", err)
	}
	c.Next()
}

func jwtSign() string {
	mySigningKey := []byte(jwtSecretKey)

	type MyCustomClaims struct {
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Date(2200, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	return "success"
}
