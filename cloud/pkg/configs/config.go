// Copyright 2020 Apulis Technology Inc. All rights reserved.

package configs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EdgeCloudConfig struct {
	DebugModel   bool
	Portal       PortalConfig
	CloudHub     CloudHubConfig
	Log          LogConfig
	Db           DbConfig
	KubeConfFile string
	KubeMaster   string
}

type HttpConfig struct {
	Address string
	Port    int
	Enable  bool
}

type PortalConfig struct {
	NodeCheckerInterval        int32
	ApplicationCheckerInterval int32
	Http                       HttpConfig
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
	Level     logrus.Level
	WriteFile bool
	FileDir   string
	FileName  string
}

type DbConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

func InitConfig(configFile string, config *EdgeCloudConfig) {
	viper.SetConfigFile(configFile)

	// set default
	viper.SetDefault("DebugModel", false)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error read config file: %s \n", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal config file: %s \n", err))
	}

	fmt.Println(config.CloudHub.Websocket)
}
