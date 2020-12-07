// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"

	"github.com/apulis/ApulisEdge/agent/pkg/app"
)

func main() {
	appUtil := app.NewInstance()

	// start server
	err := appUtil.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}
}
