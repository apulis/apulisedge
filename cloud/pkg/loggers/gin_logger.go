// Copyright 2020 Apulis Technology Inc. All rights reserved.

package loggers

import (
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinLogger(logger logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		body, _ := ioutil.ReadAll(c.Request.Body)
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"clientIP":   clientIP,
			"referer":    referer,
			"reqBody":    body,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s: %s <%d> (%dμs)", c.Request.Method, path, statusCode, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
