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

// application status
const (
	AppStatusUnpublished string = "Unpublished"
	AppStatusPublished   string = "Published"
	AppStatusOffline     string = "Offline"
)

// deploy status
const (
	StatusInit      string = "Init"
	StatusDeploying string = "Deploying"
	StatusRunning   string = "Running"
	StatusAbnormal  string = "Abnormal"
	StatusDeleting  string = "Deleting"
)

const (
	AppUserDefine string = "UserDefine"
	AppSysDefine  string = "System"
	AppTypeAll    string = "All"
)

const (
	DefaultNamespace string = "default"
)

// for ticker
const (
	TransferCountEach int = 10
)

// restart policy
const (
	RestartPolicyAlways    string = "Always"
	RestartPolicyOnFailure string = "OnFailure"
	RestartPolicyNever     string = "Never"
)

// network type
const (
	NetworkTypeHost        string = "Host"
	NetworkTypePortMapping string = "PortMapping"
)

// arch type
const (
	ArchX86 string = "x86_64"
	ArchArm string = "arm64"
)
