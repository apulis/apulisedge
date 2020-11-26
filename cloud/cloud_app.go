// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"
	"github.com/apulis/ApulisEdge/configs"
	"github.com/apulis/ApulisEdge/loggers"
	"github.com/apulis/ApulisEdge/servers/httpserver"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var logger = loggers.Log
var once sync.Once

var (
	flags      []cli.Flag
	configFile string
)

func init() {
	// init command line
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "assign config file `PATH`",
			Required:    true,
			Destination: &configFile,
		},
	}
}

func main() {
	// command line
	app := cli.NewApp()
	app.Name = "ApulisEdgeCloud"
	app.Usage = "cloud for apulis edge"
	app.Flags = flags
	app.Action = appMain

	// start server
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}

func appMain(c *cli.Context) error {
	// init config file
	initConfig()

	logger.Info("PID = %d", os.Getpid())

	// quit when signal notifys
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// start api server
	srv := httpserver.StartApiServer()

	select {
	case <-quit:
		httpserver.StopApiServer(srv)
	}

	return nil
}

func initConfig() {
	once.Do(func() {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error read config file: %s \n", err))
		}

		if err := viper.Unmarshal(&configs.CloudConfig); err != nil {
			panic(fmt.Errorf("Fatal error unmarshal config file: %s \n", err))
		}

		fmt.Println(configs.CloudConfig.CloudHub.Websocket)
	})
}
