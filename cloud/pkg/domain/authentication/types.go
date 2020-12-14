// Copyright 2020 Apulis Technology Inc. All rights reserved.

package authentication

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// interface for authentication
type Authenticator interface {
	AuthMethod(auth string) (*protocol.ApulisHeader, error)
	Init(config *configs.EdgeCloudConfig) error
}
