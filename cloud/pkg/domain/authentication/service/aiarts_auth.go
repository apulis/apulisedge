// Copyright 2020 Apulis Technology Inc. All rights reserved.

package authservice

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/apulis/ApulisEdge/cloud/pkg/protocol"
	"github.com/dgrijalva/jwt-go"
)

var logger = loggers.LogInstance()

const (
	DefaultClusterId int64 = 0
	DefaultGroupId   int64 = 0
	DefaultUserId    int64 = 0
)

type AiArtsAuthtication struct {
	JwtSecretKey  string
	Authenticator interface{}
}

type Claim struct {
	jwt.StandardClaims
	Uid      int    `json:"uid"`
	UserName string `json:"userName"`
}

func (a AiArtsAuthtication) Init(config *configs.EdgeCloudConfig) error {
	a.JwtSecretKey = config.Authentication.AiArtsAuth.Key
	logger.Debugf("jwt key = %s", a.JwtSecretKey)
	return nil
}

func (a AiArtsAuthtication) AuthMethod(auth string) (*protocol.ApulisHeader, error) {
	claim, err := a.ParseToken(auth)
	if err != nil {
		return nil, err
	}

	extracts := &protocol.ApulisHeader{
		ClusterId: DefaultClusterId,
		GroupId:   DefaultGroupId,
		UserId:    int64(claim.Uid),
	}

	return extracts, err
}

func (a AiArtsAuthtication) ParseToken(token string) (*Claim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(a.JwtSecretKey), nil
	})

	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*Claim); ok && jwtToken.Valid {
			return claim, nil
		}
	}

	logger.Debugf("AiArtsAuthtication parseToken failed! err = %v", err)
	return nil, err
}
