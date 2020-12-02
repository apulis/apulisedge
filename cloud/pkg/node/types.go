// Copyright 2020 Apulis Technology Inc. All rights reserved.

package node

import (
	_ "fmt"
)

// node status
const (
	StatusOnline       string = "Online"
	StatusOffline      string = "Offline"
	StatusNotInstalled string = "NotInstalled"
)

// node roles
const (
	EdgeRoleKey  string = "node-role.kubernetes.io/edge"
	AgentRoleKey string = "node-role.kubernetes.io/agent"
	EdgeRole     string = "edge"
	AgentRole    string = "agent"
)

const (
	TransferCountEach int = 10
)
