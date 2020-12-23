// Copyright 2020 Apulis Technology Inc. All rights reserved.

package authentication

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/domain/authentication/service"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
)

var logger = loggers.LogInstance()

var authenticatorMap = map[string]Authenticator{
	"AiArts": authservice.AiArtsAuthtication{},
}

func GetAuthenticator(config *configs.EdgeCloudConfig) (Authenticator, error) {
	authType := config.Authentication.AuthType
	_, ok := authenticatorMap[authType]
	if !ok {
		logger.Infof("InitAuth failed, authType = %s, err = %v", authType, ErrAuthTypeNotSupported)
		return nil, ErrAuthTypeNotSupported
	}

	err := authenticatorMap[authType].Init(config)
	if err != nil {
		return nil, ErrAuthenticatorInit
	}

	return authenticatorMap[authType], nil
}
