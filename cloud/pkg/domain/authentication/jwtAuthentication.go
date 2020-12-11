// Copyright 2020 Apulis Technology Inc. All rights reserved.

package authentication

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var jwtSecretKey string

type jwtAuthtication struct {
}

func (jwtAuthticator jwtAuthtication) AuthMethod(c *gin.Context) AuthResult {

	r := c.Request
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return NoAuthHeadError
	}
	tokenString := strings.Fields(auth)[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return newAuthResult(false, errors.New("jwt token parse fail"))
	}

	if token.Valid {
		return JWTAuthSuccess
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return JWTTokenFormatError
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return JWTTokenExpiredError
		} else {
			return newAuthResult(false, errors.New("can't handle this token: "+err.Error()))
		}
	} else {
		return JWTAuthFailError
	}
}

func (jwtAuthticator jwtAuthtication) initCertificate() {
	jwtSecretKey = viper.GetStringMap("authentication")["key"].(string)
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
