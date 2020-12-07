// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"
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
