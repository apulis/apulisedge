package wssclient

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kubeedge/beehive/pkg/core"
)

var wssConfig WssClientContext

func newWssClient(enable bool) *WssClientContext {
	return &WssClientContext{}
}

// Register registers the module
func Register(wss *WssClientContext) {
	logger = wss.Logger
	core.Register(newWssClient(true))
}

// Name get name of the module
func (wss *WssClientContext) Name() string {
	return "WssClientConfig"
}

// Group get group of the module
func (wss *WssClientContext) Group() string {
	return "WssClientConfig"
}

// Enable indicates whether the module is enabled
func (wss *WssClientContext) Enable() bool {
	return true
}

// Start runs the module
func (wss *WssClientContext) Start() {
	for {

		u := url.URL{
			Scheme:      "ws",
			Opaque:      "",
			User:        &url.Userinfo{},
			Host:        wssConfig.Server + ":" + strconv.Itoa(wssConfig.Port),
			Path:        "",
			RawPath:     "",
			ForceQuery:  false,
			RawQuery:    "",
			Fragment:    "",
			RawFragment: "",
		}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logger.Errorf("Can't access server ")
			continue
		}
		defer c.Close()
		_, msg, err := c.ReadMessage()
		if err != nil {
			logger.Errorf("Can't read message")
			continue
		}
		msgString := string(msg)
		logger.Infof("Recieve message: " + msgString)

		time.Sleep(time.Duration(3) * time.Second)
	}
}
