// Copyright 2020 Apulis Technology Inc. All rights reserved.

package configs

type EdgeCloudConfig struct {
	Portal   PortalConfig
	CloudHub CloudHubConfig
	Log      LogConfig
}

type HttpConfig struct {
	Address string
	Port    int
	Enable  bool
}

type PortalConfig struct {
	Http HttpConfig
}

type WebsocketConfig struct {
	Address string
	Port    int
	Enable  bool
}

type CloudHubConfig struct {
	Websocket WebsocketConfig
}

type LogConfig struct {
	WriteFile bool
	FileDir   string
	FileName  string
}

var CloudConfig EdgeCloudConfig
