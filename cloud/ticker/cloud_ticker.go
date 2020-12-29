// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"os"
)

var logger = loggers.LogInstance()

func main() {
	// create cloud app and initialize
	app := CloudTickerInstance()
	if err := app.Init("ApulisEdgeTicker", "Handle asynchronous task from cloud, fetch the task from db or messageQueue"); err != nil {
		logger.Fatalf("Application Initialize Failed. Error = %v", err)
		os.Exit(1)
	}

	// start cloud ticker
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
