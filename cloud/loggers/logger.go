// Copyright 2020 Apulis Technology Inc. All rights reserved.

package loggers

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/apulis/ApulisEdge/configs"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	logConf := configs.CloudConfig.Log
	Log.Formatter = new(Formatter)

	if logConf.WriteFile {
		if err := os.Mkdir(logConf.FileDir, 0755); err != nil {
			fmt.Println(err.Error())
		}
		fileName := path.Join(logConf.FileDir, logConf.FileName)
		if _, err := os.Stat(fileName); err != nil {
			if _, err := os.Create(fileName); err != nil {
				fmt.Println(err.Error())
			}
		}
		writeToFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
		}

		writers := io.MultiWriter(os.Stdout, writeToFile)
		Log.Out = writers
	}
}
