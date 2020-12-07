// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"
	"github.com/apulis/ApulisEdge/agent/pkg/app"
	_ "github.com/apulis/ApulisEdge/agent/pkg/model_service/service"
	"os"
)

func main() {
	appUtil := app.NewInstance()
	client.Run()

	// start server
	err := appUtil.Run(os.Args)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
