package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of cobrademo",
		Long:  `All software has versions. This is cobrademo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cobrademo version is v1.0")
		},
	}
	cmd.AddCommand(versionCmd)
	initApp()

	return cmd
}

// Run is the start function
func run() error {
	initConfig()
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
	viper.SetConfigFile("/etc/apulisedge/config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s", err))
	}

}

func initApp() {

}
