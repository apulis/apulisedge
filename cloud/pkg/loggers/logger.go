// Copyright 2020 Apulis Technology Inc. All rights reserved.

package loggers

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var instance *logrus.Logger

var logger = LogInstance()

func LogInstance() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
	})
	return instance
}

func InitLogger(config *configs.EdgeCloudConfig) {
	logConf := config.Log
	logger.Formatter = new(Formatter)
	logger.SetLevel(config.Log.Level)

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
		logger.Out = writers
	}
}
