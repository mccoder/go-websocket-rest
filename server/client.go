package server

import (
	"fmt"
	"io"
	"log"
	"net"

	wsm "github.com/mccoder/go-websocket-rest/model"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
)

const channelBufSize = 100
const err_ActionNotFound = "Action not found"

type WsClientOp func(c *WsClientListener)

// Chat client.
type WsClientListener struct {
	SessionId  uint64
	User       WsUser
	RemoteAddr string
	wsClient   *websocket.Conn
	server     *WsServer

	chOP          chan WsClientOp
	doneCh        chan bool
	subscriptions map[string]string
}

// Create new chat client.
func newWsClient(user WsUser, ws *websocket.Conn, server *WsServer, RemoteAddr net.Addr) *WsClientListener {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}
	doneCh := make(chan bool, 2)
	chOP := make(chan WsClientOp, channelBufSize)
	subscriptions := make(map[string]string)
	return &WsClientListener{
		User:          user,
		wsClient:      ws,
		RemoteAddr:    RemoteAddr.String(),
		server:        server,
		chOP:          chOP,
		doneCh:        doneCh,
		subscriptions: subscriptions,
	}
}

//func (c *Client) Conn() *websocket.Conn {
//	return c.wsClient
//}

func WsWithWriteJsonMsg(msg *wsm.WsResponceMessage) WsClientOp {
	return func(c *WsClientListener) {
		c.wsClient.ReadJSON(msg)
	}
}

/**
* Отправка уведомлений.
 */
func (c *WsClientListener) Notify(topicId string, msg *wsm.WsResponceMessage) {
	id, ok := c.subscriptions[topicId]
	if ok {
		msg.Id = id
		msg.Type = wsm.TypeTopic
		select {
		case c.chOP <- WsWithWriteJsonMsg(msg):
		default:
			c.server.Del(c)
			err := fmt.Errorf("client %d is disconnected.", c.SessionId)
			c.server.Err(err)
		}
	}
}

func (c *WsClientListener) Write(msg *wsm.WsResponceMessage) {
	select {
	case c.chOP <- WsWithWriteJsonMsg(msg):
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.SessionId)
		c.server.Err(err)
	}
}

func (c *WsClientListener) Done() {
	c.doneCh <- true // for listenWrite method
	c.doneCh <- true // for listenRead method
}

// Listen Write and Read request via chanel
func (c *WsClientListener) Listen() {
	go c.listenWrite(c.chOP)
	c.listenRead()
}

// Listen write request via chanel
func (c *WsClientListener) listenWrite(chOP <-chan WsClientOp) {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case op := <-chOP:
			//log.Println("Send:", msg)
			op(c)

		// receive done request
		case <-c.doneCh:
			for name, _ := range c.subscriptions {
				c.server.Unsubscribe(c, name)
			}

			c.server.Del(c)

			//c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *WsClientListener) processSubscribeMessage(msg *wsm.WsRequestMessage, response *wsm.WsResponceMessage) {
	request := &wsm.WsSubscribeAction{}
	err := easyjson.Unmarshal(msg.ParamsByteArray(), request)
	if err != nil {
		response.Code = 501
		response.Error = "Invalid susbscribe request"
	} else {
		response.Code = 410
		_, found := c.subscriptions[request.TopicId]
		if !found {
			c.subscriptions[request.TopicId] = response.Id
			c.server.Subscribe(c, request.TopicId)
		}
	}
}

func (c *WsClientListener) processUnsubscribeMessage(msg *wsm.WsRequestMessage, response *wsm.WsResponceMessage) {
	request := &wsm.WsSubscribeAction{}
	err := easyjson.Unmarshal(msg.ParamsByteArray(), request)
	if err != nil {
		response.Code = 501
		response.Error = "Invalid unsusbscribe request"
	} else {
		response.Code = 410
		_, found := c.subscriptions[request.TopicId]
		if found {
			delete(c.subscriptions, request.TopicId)
			c.server.Unsubscribe(c, request.TopicId)
		}
	}
}

func (c *WsClientListener) processMessage(msg *wsm.WsRequestMessage, response *wsm.WsResponceMessage) {
	switch msg.Action {
	case "$ping":
		response.Code = 200

	case "$subscribe":
		c.processSubscribeMessage(msg, response)

	case "$unsubscribe":
		c.processUnsubscribeMessage(msg, response)

	default:
		log.Printf("WS action %s\n", msg.Action)

		callback, existsCallback := c.server.callbacks[msg.Action]
		if !existsCallback {
			response.Error = err_ActionNotFound
			response.Code = 404
			log.Printf("WS action '%s' not found\n", msg.Action)
		} else {
			res_data, err := callback(msg.Params, c)
			if err != nil {
				response.Error = err.Error()
				response.Code = 501
			} else {
				response.Data = res_data
				response.Code = 200
			}
		}
	}
}

// Listen read request via chanel
func (c *WsClientListener) listenRead() {
	log.Println("Listening read from client")
	msg := &wsm.WsRequestMessage{}
	response := &wsm.WsResponceMessage{Code: 500}
	for {
		select {

		// receive done request
		case <-c.doneCh:
			//c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			{
				msg.Clean()
				err := c.wsClient.ReadJSON(&msg)
				if err == nil {
					response.Clear()
					response.Id = msg.Id
					c.processMessage(msg, response)

					err := c.wsClient.WriteJSON(response)
					if err != nil {
						log.Println("Cant serialize responce")
					}
					response.Clear()
					msg.Clean()
				} else {
					//пока в случае ошибок, принудительно разрваем соединение с клиентом.
					//Возможно он начал отправлять сообщения в не известном формате.
					if err == io.EOF {
						c.Done()
					} else {
						log.Printf("WS client err: %s\n", err.Error())
						c.server.Err(err)
						c.Done()
					}
				}
			}
		}
	}
}
