package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apulis/ApulisEdge/agent/pkg/common/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

	initApp()

	return cmd
}

// Run is the start function
func run() error {

	initConfig()

	initLogger()

	fmt.Println("app running")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			fmt.Println("still running")
		}
	}()

	fmt.Println(viper.Get("server"))
	select {
	case <-quit:
		fmt.Printf("app quit")
	}
	return nil
}

func registerModules() error {
	return nil
}

func initConfig() {
	config.InitConfig()
}

func initApp() {

}

func initLogger() {

}
