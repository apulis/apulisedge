package config

import (
	"fmt"

	"github.com/apulis/ApulisEdge/agent/pkg/common/loggers"
	wssclient "github.com/apulis/ApulisEdge/agent/pkg/wssClient"
	"github.com/spf13/viper"
)

// ConfigFilePath indicates runtime configuration file path
var ConfigFilePath string

// AppConfig indicates runtime configuration
var AppConfig AgentConfigType

// AgentConfigType indicates configuration structure
type AgentConfigType struct {
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Modules  Modules
}

// LogConfig indicates log configuration structure
type LogConfig struct {
	Path     string `yaml:"path"`
	FileName string `yaml:"filename"`
}

type DatabaseConfig struct {
	Dir      string `yaml:"dir"`
	Type     string `yaml:"type"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Modules struct {
	Logger    *loggers.LoggerConfigType
	WssClient *wssclient.WssClientContext
}

func InitConfig() {
	viper.SetConfigFile(ConfigFilePath)

	viper.SetDefault("database.databaseDir", DEFAULT_DATABASE_DIR)
	viper.SetDefault("log.path", DEFAULT_LOG_PATH)
	viper.SetDefault("log.filename", DEFAULT_LOG_FILE)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s", err))
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(fmt.Errorf("Fatal error convert config file: %s", err))
	}
	moduleInit()
}

func moduleInit() {
	AppConfig.Modules.Logger = &loggers.LoggerConfigType{
		Path:     AppConfig.Log.Path,
		FileName: AppConfig.Log.FileName,
	}

	AppConfig.Modules.WssClient = &wssclient.WssClientContext{
		Logger:   loggers.LogInstance(),
		IsEnable: true,
		Server:   AppConfig.Server.Address,
		Port:     AppConfig.Server.Port,
	}

}
