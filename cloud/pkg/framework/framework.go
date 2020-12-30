// Copyright 2020 Apulis Technology Inc. All rights reserved.

package framework

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/database"
	appentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/application/entity"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"github.com/apulis/ApulisEdge/cloud/pkg/servers/httpserver"
	"github.com/urfave/cli/v2"
)

var logger = loggers.LogInstance()

type CloudApp struct {
	internalApp      *cli.App
	flags            []cli.Flag
	configFile       string
	cloudConfig      configs.EdgeCloudConfig
	tickerCancelFunc context.CancelFunc
	tickerCtx        context.Context
	clusters         []configs.ClusterConfig
}

var once sync.Once
var instance *CloudApp

func CloudAppInstance() *CloudApp {
	once.Do(func() {
		instance = &CloudApp{
			internalApp: cli.NewApp(),
		}
	})
	return instance
}

func (app *CloudApp) Init(appName string, appUsage string) error {
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

func (app *CloudApp) Run(arguments []string) error {
	err := app.internalApp.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func (app *CloudApp) MainLoop() error {
	logger.Infof("PID = %d", os.Getpid())

	// init config
	app.InitConfig()

	// init logger
	app.InitLogger()

	// init database
	app.InitDatabase()

	// init tables
	app.InitTables()

	// init auth
	err := app.InitAuth()
	if err != nil {
		logger.Errorf("Auth init failed, err = %v", err)
		return err
	}

	// init clusters
	app.InitClusters()

	// quit when signal notifys
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// start api server
	srv := httpserver.StartApiServer(&app.cloudConfig)

	select {
	case <-quit:
		app.tickerCancelFunc()
		httpserver.StopApiServer(srv)
	}

	return nil
}

func (app *CloudApp) InitLogger() {
	loggers.InitLogger(&app.cloudConfig)
}

func (app *CloudApp) InitConfig() {
	logger.Infof("InitConfig, configFile path = %s", app.configFile)
	configs.InitConfig(app.configFile, &app.cloudConfig)
}

func (app *CloudApp) InitDatabase() {
	database.InitDatabase(&app.cloudConfig)
}

func (app *CloudApp) InitAuth() error {
	return httpserver.InitAuth(&app.cloudConfig)
}

func (app *CloudApp) InitClusters() {
	cluster.InitClusters(&app.cloudConfig)
}

func (app *CloudApp) InitTables() {
	database.CreateTableIfNotExists(nodeentity.NodeBasicInfo{})
	database.CreateTableIfNotExists(nodeentity.NodeOfBatchInfo{})
	database.CreateTableIfNotExists(nodeentity.BatchTaskRecord{})
	database.CreateTableIfNotExists(appentity.ApplicationBasicInfo{})
	database.CreateTableIfNotExists(appentity.ApplicationVersionInfo{})
	database.CreateTableIfNotExists(appentity.ApplicationDeployInfo{})
	database.CreateTableIfNotExists(imageentity.UserContainerImageInfo{})
	database.CreateTableIfNotExists(imageentity.UserContainerImageVersionInfo{})
	database.CreateTableIfNotExists(imageentity.ContainerImageOrg{})
}
