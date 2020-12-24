// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	"errors"
)

var (
	ErrNodeNotExist         = errors.New("node not exist")
	ErrNodeTypeNotExist     = errors.New("node type not exist")
	ErrGetNode              = errors.New("get node err")
	ErrDeleteStatusDeleting = errors.New("delete failed! status is deleting")
)
