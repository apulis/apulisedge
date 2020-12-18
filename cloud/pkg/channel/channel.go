// Copyright 2020 Apulis Technology Inc. All rights reserved.

package channel

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"sync"
)

var once sync.Once
var instance *ChanContext

var logger = loggers.LogInstance()

// constants for channel context
const (
	ChanSizeDefault = 1024
)

const (
	ModuleNameContainerImage = "containerImage"
)

type ChanMsgType chan interface{}

type ChanContext struct {
	channels map[string]ChanMsgType
	chsLock  sync.RWMutex
}

func ChanContextInstance() *ChanContext {
	once.Do(func() {
		channelMap := make(map[string]ChanMsgType)
		instance = &ChanContext{
			channels: channelMap,
			chsLock:  sync.RWMutex{},
		}
	})
	return instance
}

func (ctx *ChanContext) NewChannel() ChanMsgType {
	channel := make(ChanMsgType, ChanSizeDefault)
	return channel
}

func (ctx *ChanContext) AddChannel(module string, moduleCh ChanMsgType) {
	ctx.chsLock.Lock()
	defer ctx.chsLock.Unlock()

	ctx.channels[module] = moduleCh
}

// deleteChannel by module name
func (ctx *ChanContext) DelChannel(module string) {
	// delete module channel from channels map
	ctx.chsLock.Lock()
	defer ctx.chsLock.Unlock()

	_, exist := ctx.channels[module]
	if !exist {
		logger.Warningf("Failed to get channel, module:%s", module)
		return
	}
	delete(ctx.channels, module)
}

// getChannel return chan
func (ctx *ChanContext) GetChannel(module string) ChanMsgType {
	ctx.chsLock.RLock()
	defer ctx.chsLock.RUnlock()

	if _, exist := ctx.channels[module]; exist {
		return ctx.channels[module]
	}

	logger.Warningf("Failed to get channel, type:%s", module)
	return nil
}
