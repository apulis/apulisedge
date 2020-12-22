// Copyright 2020 Apulis Technology Inc. All rights reserved.

package domain

import (
	"errors"
	_ "fmt"
)

// arch type
const (
	ArchX86 string = "x86_64"
	ArchArm string = "arm64"
)

var (
	ErrArchTypeNotExist = errors.New("arch type not exist")
)

func IsArchValid(arch string) bool {
	if arch != ArchX86 && arch != ArchArm {
		return false
	}

	return true
}
