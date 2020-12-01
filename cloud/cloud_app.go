// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"github.com/apulis/ApulisEdge/cloud/framework"
	"github.com/apulis/ApulisEdge/cloud/loggers"
	"os"
)

var logger = loggers.LogInstance()

func main() {
	// create cloud app and initialize
	app := framework.CloudAppInstance()
	if err := app.Init("ApulisEdgeCloud", "ApulisEdgeCloud"); err != nil {
		logger.Fatalf("Application Initialize Failed. Error = %v", err)
		os.Exit(1)
	}

	// start cloud app
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
