package wssclient

import (
	"github.com/sirupsen/logrus"
)

type WssClientContext struct {
	Logger   *logrus.Logger
	IsEnable bool
	Server   string
	Port     int
}

var logger *logrus.Logger
