// Copyright 2020 Apulis Technology Inc. All rights reserved.

package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

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
	app.Name = "ApulisEdgeAgent"
	app.Usage = "agent for apulis edge"
	app.Flags = flags
	app.Action = appMain

	// start server
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func appMain(c *cli.Context) error {
	return nil
}
