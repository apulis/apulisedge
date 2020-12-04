// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	"errors"
)

var (
	ErrDeployExist    = errors.New("application deploy exist")
	ErrGetApplication = errors.New("get application err")
)
