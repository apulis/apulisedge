package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apulis/ApulisEdge/agent/pkg/agentSocket"
	"github.com/apulis/ApulisEdge/agent/pkg/common/config"
	"github.com/apulis/ApulisEdge/agent/pkg/common/database"
	"github.com/apulis/ApulisEdge/agent/pkg/common/loggers"
	"github.com/spf13/cobra"
)

var logger = loggers.LogInstance()

// NewAgentCommand create a new instance of the application
func NewAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Short: "agent for apulisedge, use with kubeedge",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	cmd.Flags().StringVarP(&config.ConfigFilePath, "config", "c", config.DEFAULT_CONFIG_PATH, "path to config file")

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Print the version number of cobrademo",
		Long:  `All software has versions. This is cobrademo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.ConfigFilePath)
		},
	}
	testCmd.Flags().StringVarP(&config.ConfigFilePath, "config", "c", "/etc/apulisedge/config/config.yaml", "path to config file")
	cmd.AddCommand(testCmd)

	return cmd
}

// Run is the start function
func run() error {

	initConfig()

	initLogger()

	initDatabase()

	logger.Infoln("app start, config showing bellow:")
	logger.Infoln("============================== ")
	logger.Infoln(config.AppConfig)
	logger.Infoln("============================== ")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			logger.Infoln("app still running ...")
		}
	}()

	select {
	case <-quit:
		fmt.Printf("app quit")
	}

	stopApp()
	return nil
}

func registerModules() error {
	return nil
}

func initConfig() {
	config.InitConfig()
}

func initLogger() {
	loggers.InitLogger()

}

func initDatabase() {
	database.InitDatabase()

	database.CreateTableIfNotExists(agentSocket.WebSocket{})
}

func stopApp() {
	database.CloseDatabase()
}
