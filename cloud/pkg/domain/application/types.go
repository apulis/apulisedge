// Copyright 2020 Apulis Technology Inc. All rights reserved.

package application

import (
	_ "fmt"
)

// application status
const (
	StatusUnPublished string = "UnPublished"
	StatusPublished   string = "Published"
)

// deploy status
const (
	StatusInit      string = "Init"
	StatusDeploying string = "Deploying"
	StatusRunning   string = "Running"
	StatusDeleting  string = "Deleting"
)

const (
	AppUserDefine string = "UserDefine"
	AppSysDefine  string = "System"
)

const (
	DefaultNamespace string = "default"
)

// for ticker
const (
	TransferCountEach int = 10
)
