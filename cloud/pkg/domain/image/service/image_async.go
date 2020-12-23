// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/channel"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"time"
)

// image async handler
func ImageAsyncHandleLoop(ctx context.Context, config *configs.EdgeCloudConfig) {
	duration := time.Duration(config.ContainerImage.ImageCheckerInterval) * time.Second
	checkTicker := time.NewTimer(duration)
	defer checkTicker.Stop()

	msgChanContext := channel.ChanContextInstance()
	msgChan := msgChanContext.GetChannel(channel.ModuleNameContainerImage)

	for {
		select {
		case <-ctx.Done():
			logger.Infof("ImageAsyncHandleLoop was terminated")
			return
		case msg := <-msgChan:
			ImageAsyncHandle(msg)
		}
	}
}

func ImageAsyncHandle(msg interface{}) {
	switch msg.(type) {
	case imageentity.UserContainerImageVersionInfo:
		logger.Infof("handle msg delete image = %+v", msg)
		// TODO del image from harbor
	default:
		logger.Infof("Not Support msg = %+v", msg)
	}
}
