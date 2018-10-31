package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WsNewClientHandler func(client *WsClientListener)
type WsClientSubscribedHandler func(channelId string, client *WsClientListener)

type WsHandlerCallback func(data string, client *WsClientListener) (interface{}, error)

type WsHandler interface {
	WsPrefix() string
	GetCallback() map[string]WsHandlerCallback
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func WsWithUserHandler(w http.ResponseWriter, r *http.Request, user WsUser, wsServer *WsServer) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsServer.NewClientConnected(user, conn)
}
