// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	"errors"
)

var (
	ErrDeployExist              = errors.New("application deploy exist")
	ErrApplicationExist         = errors.New("application exist")
	ErrApplicationVersionExist  = errors.New("application version exist")
	ErrImageVersionNotExist     = errors.New("image version not exist")
	ErrUnDeploying              = errors.New("undeploying")
	ErrChangeAppVersionFailed   = errors.New("change app version fail")
	ErrDeleteStatusNotOffline   = errors.New("delete failed! status not offline")
	ErrDeleteStatusPublished    = errors.New("delete failed! status is published")
	ErrDeployStatusNotPublished = errors.New("deploy failed! status not published")
	ErrWrongRestartPolicy       = errors.New("wrong restart policy")
	ErrNetworkPortmappingEmpty  = errors.New("need port mapping")
	ErrDeployPartFails          = errors.New("part of deploy fails")
	ErrDeployAllFails           = errors.New("all of deploy fails")
	ErrUnDeployPartFails        = errors.New("part of undeploy fails")
	ErrUnDeployAllFails         = errors.New("all of undeploy fails")
	ErrArchTypeNotExist         = errors.New("arch type not exist")
)
