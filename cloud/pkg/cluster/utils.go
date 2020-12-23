// Copyright 2020 Apulis Technology Inc. All rights reserved.

package cluster

import (
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
)

func (c *Cluster) GetUniqueName(resourceName string) string {
	return resourceName + "-" + UUIDToString(uuid.NewV4())
}

func IsArchValid(arch string) bool {
	if arch != ArchX86 && arch != ArchArm {
		return false
	}

	return true
}

func UUIDToString(id uuid.UUID) string {
	buf := make([]byte, 32)
	hex.Encode(buf[0:8], id[0:4])
	hex.Encode(buf[8:12], id[4:6])
	hex.Encode(buf[12:16], id[6:8])
	hex.Encode(buf[16:20], id[8:10])
	hex.Encode(buf[20:], id[10:])
	return string(buf)
}
