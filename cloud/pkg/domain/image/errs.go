// Copyright 2020 Apulis Technology Inc. All rights reserved.

package image

import (
	"errors"
)

var (
	ErrImageVersionExist = errors.New("image version exist")
	ErrOrgImageNotEmpty  = errors.New("org images not empty")
)
