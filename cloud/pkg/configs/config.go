// Copyright 2020 Apulis Technology Inc. All rights reserved.

package configs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EdgeCloudConfig struct {
	DebugModel     bool
	Portal         PortalConfig
	CloudHub       CloudHubConfig
	Log            LogConfig
	Db             DbConfig
	Authentication AuthConfig
	Clusters       []ClusterConfig
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

type AuthConfig struct {
	AuthType   string
	AiArtsAuth AiArtsAuthConfig
	// TODO add other auth
}

type AiArtsAuthConfig struct {
	Key string
}

type ClusterConfig struct {
	Id            int64
	Desc          string
	KubeMaster    string
	KubeConfFile  string
	HarborUser    string
	HarborPasswd  string
	HarborAddress string
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

	fmt.Printf("Portal config = %v\n", config.Portal)
	fmt.Printf("Websocket config = %v\n", config.CloudHub.Websocket)
	fmt.Printf("Cluster config = %v\n", config.Clusters)
}
