// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	"errors"
)

var (
	ErrDeployExist             = errors.New("application deploy exist")
	ErrApplicationExist        = errors.New("application exist")
	ErrApplicationVersionExist = errors.New("application version exist")
	ErrUnDeploying             = errors.New("undeploying")
)
