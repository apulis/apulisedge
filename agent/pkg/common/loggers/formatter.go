package loggers

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type logFormatter struct {
}

func (formatter *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil

}
