package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ConfigFilePath indicates runtime configuration file path
var ConfigFilePath string

// AppConfig indicates runtime configuration
var AppConfig AgentConfig

// AgentConfig indicates configuration structure
type AgentConfig struct {
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// LogConfig indicates log configuration structure
type LogConfig struct {
	Path     string `yaml:"path"`
	FileName string `yaml:"filename"`
}

type DatabaseConfig struct {
	DatabaseDir     string `yaml:"dir"`
	DatabaseType    string `yaml:"type"`
	DatabaseAddress string `yaml:"address"`
	Port            int    `yaml:"port"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
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

	fmt.Println("============================== ")
	fmt.Println(AppConfig)
	fmt.Println("============================== ")
}
