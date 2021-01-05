// Copyright 2020 Apulis Technology Inc. All rights reserved.

package configs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sigs.k8s.io/yaml"
)

type EdgeCloudConfig struct {
	DebugModel     bool
	Portal         PortalConfig
	ContainerImage ImageConfig
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

type ImageConfig struct {
	ImageCheckerInterval int32
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
	Id              int64
	Desc            string
	Domain          string
	KubeMaster      string
	KubeConfFile    string
	HarborAddress   string
	HarborProject   string
	HarborUser      string
	HarborPasswd    string
	DownloadAddress string
}

func InitConfig(configFile string, config *EdgeCloudConfig) {
	viper.SetConfigFile(configFile)

	// set default
	viper.SetDefault("DebugModel", false)
	viper.SetDefault("ContainerImage.ImageCheckerInterval", 10)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error read config file: %s \n", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal config file: %s \n", err))
	}

	fmt.Printf("Portal config = %v\n", config.Portal)
	fmt.Printf("Image config = %v\n", config.ContainerImage)
	fmt.Printf("Websocket config = %v\n", config.CloudHub.Websocket)
	fmt.Printf("Cluster config = %+v\n", config.Clusters)
}

func PrintMinConfigAndExitIfRequested(config interface{}) {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("Marshal min config to yaml error %v", err)
		os.Exit(1)
	}
	fmt.Println("# With --minconfig , you can easily used this configurations as reference.")
	fmt.Println("# It's useful to users who are new to ApulisEdge, and you can modify/create your own configs accordingly. ")
	fmt.Println("# This configuration is suitable for beginners.")
	fmt.Printf("\n%v\n\n", string(data))
	os.Exit(0)
}

// NewMinCloudCoreConfig returns a min ClusterConfig object
func NewMinCloudConfig() *EdgeCloudConfig {
	return &EdgeCloudConfig{
		Portal: PortalConfig{
			NodeCheckerInterval:        30,
			ApplicationCheckerInterval: 30,
			Http: HttpConfig{
				Enable:  true,
				Address: "0.0.0.0",
				Port:    32769,
			},
		},
		CloudHub: CloudHubConfig{
			Websocket: WebsocketConfig{
				Enable:  false,
				Address: "0.0.0.0",
				Port:    32768,
			},
		},
		Log: LogConfig{
			Level:     4,
			WriteFile: true,
			FileDir:   "/var/log/apulisedge",
			FileName:  "apulis_edge_cloud.log",
		},
		Db: DbConfig{
			Username:     "user",
			Password:     "password",
			Host:         "127.0.0.1",
			Port:         3306,
			Database:     "ApulisEdgeCloudDB",
			MaxIdleConns: 2,
			MaxOpenConns: 10,
		},
		Authentication: AuthConfig{
			AuthType: "AiArts",
			AiArtsAuth: AiArtsAuthConfig{
				Key: "jwt sign key",
			},
		},
		Clusters: []ClusterConfig{
			{
				Id:              0,
				Desc:            "for test",
				Domain:          "edge.yourcorp.com",
				KubeMaster:      "https://c0.edge.yourcorp.com",
				KubeConfFile:    "/root/.kube/config",
				HarborAddress:   "harbor.yourcorp.com",
				HarborProject:   "apulisedge",
				HarborUser:      "user",
				HarborPasswd:    "password",
				DownloadAddress: "https://c0.download.yourcorp.com",
			},
		},
	}
}
