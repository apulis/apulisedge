package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Uid      int    `json:"uid"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

var jwtSecretKey = []byte("Sign key for JWT")

func GenerateToken(uid int, userName string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		uid,
		userName,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecretKey)

	fmt.Printf("generated token = %s\n", token)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecretKey, nil
	})

	fmt.Printf("jwtToken = %v, claim = %v\n", jwtToken, jwtToken.Claims.(*Claims))
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
			return claim, nil
		}
	}

	fmt.Printf("AiArtsAuthtication parseToken failed! err = %v\n", err)
	return nil, err
}

func main() {
	token, err := GenerateToken(30001, "admin")
	if err != nil {
		fmt.Printf("generate token failed! err = %v\n", err)
		return
	}

	claims, err := ParseToken(token)
	if err != nil {
		fmt.Printf("parse token failed! err = %v\n", err)
		return
	}

	fmt.Printf("parse token succ!, cliams = %v\n", claims)
}
