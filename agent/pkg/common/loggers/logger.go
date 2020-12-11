package loggers

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/apulis/ApulisEdge/agent/pkg/common/config"
	"github.com/sirupsen/logrus"
)

var instance *logrus.Logger
var once sync.Once

var logger = LogInstance()

// InitLogger initializes log module
func InitLogger() {
	// create log dir
	_, err := os.Stat(config.AppConfig.Log.Path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(config.AppConfig.Log.Path, 0755)
		} else {
			logger.Panicln("Fatal log directory error: %s", err)
		}
	}
	// create log file
	logFilePath := path.Join(config.AppConfig.Log.Path, config.AppConfig.Log.FileName)
	fileWriter, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Panicln("Fatal error open log file: %s", err)
	}
	// set log output
	mw := io.MultiWriter(os.Stdout, fileWriter)
	logger.SetOutput(mw)
	// set log format
	logger.SetFormatter(new(logFormatter))
}

// LogInstance acquire app logger
func LogInstance() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
	})
	return instance

}
