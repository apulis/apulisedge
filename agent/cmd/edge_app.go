// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"

	"github.com/apulis/ApulisEdge/agent/pkg/app"
)

func main() {
	command := app.NewAgentCommand()

	// start server
	err := command.Execute()
	if err != nil {
		fmt.Printf(err.Error())
	}
}
