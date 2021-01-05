// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apulis/ApulisEdge/cloud/pkg/channel"
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	"github.com/apulis/ApulisEdge/cloud/pkg/database"
	applicationticker "github.com/apulis/ApulisEdge/cloud/ticker/domain/application"
	"github.com/apulis/ApulisEdge/cloud/ticker/domain/batchinstall"
	nodeticker "github.com/apulis/ApulisEdge/cloud/ticker/domain/node"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/urfave/cli/v2"
)

type CloudTicker struct {
	internalApp      *cli.App
	flags            []cli.Flag
	configFile       string
	cloudConfig      configs.EdgeCloudConfig
	tickerCancelFunc context.CancelFunc
	tickerCtx        context.Context
	clusters         []configs.ClusterConfig
}

var once sync.Once
var instance *CloudTicker

func CloudTickerInstance() *CloudTicker {
	once.Do(func() {
		instance = &CloudTicker{
			internalApp: cli.NewApp(),
		}
	})
	return instance
}

func (app *CloudTicker) Init(appName string, appUsage string) error {
	app.flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "assign config file `PATH`",
			Value:       "/etc/apulisedge/cloud/cloud.yaml",
			EnvVars:     []string{"APULIS_EDGE_CLOUD_CONFIG"},
			Destination: &app.configFile,
		},
	}

	app.internalApp.Name = appName
	app.internalApp.Usage = appUsage
	app.internalApp.Flags = app.flags
	app.internalApp.Action = func(c *cli.Context) error {
		return app.MainLoop()
	}

	app.tickerCtx, app.tickerCancelFunc = context.WithCancel(context.Background())
	return nil
}

func (app *CloudTicker) Run(arguments []string) error {
	err := app.internalApp.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func (app *CloudTicker) MainLoop() error {
	logger.Infof("PID = %d", os.Getpid())

	// init config
	app.InitConfig()

	// init logger
	app.InitLogger()

	// init database
	app.InitDatabase()

	// init clusters
	app.InitClusters()

	// msg chan
	msgChanContext := channel.ChanContextInstance()
	msgChanContext.AddChannel(channel.ModuleNameContainerImage, msgChanContext.NewChannel())

	// init ticker
	go nodeticker.CreateNodeTickerLoop(app.tickerCtx, &app.cloudConfig)
	go applicationticker.CreateApplicationTickerLoop(app.tickerCtx, &app.cloudConfig)
	go batchinstall.CreateBatchInstallTicker(app.tickerCtx, &app.cloudConfig)

	// quit when signal notifys
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		app.tickerCancelFunc()
	}

	return nil
}

func (app *CloudTicker) InitLogger() {
	loggers.InitLogger(&app.cloudConfig)
}

func (app *CloudTicker) InitConfig() {
	logger.Infof("InitConfig, configFile path = %s", app.configFile)
	configs.InitConfig(app.configFile, &app.cloudConfig)
}

func (app *CloudTicker) InitDatabase() {
	database.InitDatabase(&app.cloudConfig)
}

func (app *CloudTicker) InitClusters() {
	cluster.InitClusters(&app.cloudConfig)
}
