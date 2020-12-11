package modelservice

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func Run() {
	addr := "127.0.0.1"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	websocket.DefaultDialer.Dial(u.String(), nil)

}
