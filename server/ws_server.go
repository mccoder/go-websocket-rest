package server

import (
	"github.com/gorilla/websocket"
	. "github.com/mccoder/go-websocket-rest/model"
	"log"
	"math"
	"runtime"
)

type WsUser interface {
	GetId() uint64
}

type subscribeEvent struct {
	client *WsClientListener
	event  string
}

type WsServerOperation func(w *WsServer)

type WsServer struct {
	nextSessionId uint64
	clientMap     map[uint64]*WsClientListener

	chanOp chan WsServerOperation
	doneCh chan bool

	errCh chan error
	//Признак того, что сервер работает.
	//Выставляется перед остановкой, чтобы новые клиенты больше не регистрировались.
	working bool

	callbacks               map[string]WsHandlerCallback
	subscribeCallbacks      map[string][]*WsClientListener
	NewClientHandler        WsNewClientHandler
	SubscribeClientHandler  WsClientSubscribedHandler
	subscribedClientHandler map[string][]WsNewClientHandler
}

func deleteClient(clients []*WsClientListener, c *WsClientListener) []*WsClientListener {
	for i := range clients {
		if clients[i] == c {
			return append(clients[:i], clients[i+1:]...)
			break
		}
	}
	return clients
}

// Create new server.
func NewWsServer() *WsServer {
	callbacks := make(map[string]WsHandlerCallback)
	subscribeCallbacks := make(map[string][]*WsClientListener)
	chanOp := make(chan WsServerOperation)

	doneCh := make(chan bool, 1)
	errCh := make(chan error)

	server := &WsServer{
		working:                 true,
		chanOp:                  chanOp,
		doneCh:                  doneCh,
		errCh:                   errCh,
		callbacks:               callbacks,
		subscribeCallbacks:      subscribeCallbacks,
		clientMap:               make(map[uint64]*WsClientListener),
		subscribedClientHandler: make(map[string][]WsNewClientHandler),
	}

	runtime.SetFinalizer(server, func(v interface{}) {
		i, ok := v.(*WsServer)
		if !ok {
			log.Printf("finalizer called with type %T, want *bigValue", v)
		}

		i.doneClients()
		i.closeChannels()
	})

	return server
}

func (s *WsServer) closeChannels() {
	close(s.chanOp)
	close(s.doneCh)
	close(s.errCh)
}

func (s *WsServer) findNextSessionId() uint64 {
	if s.nextSessionId == math.MaxUint64 {
		s.nextSessionId = 0
	}
	sessionId := s.nextSessionId + 1
	for {
		if _, found := s.clientMap[sessionId]; !found {
			s.nextSessionId = sessionId
			break
		}
		if sessionId == math.MaxUint64 {
			sessionId = 0
		} else {
			sessionId++
		}
	}
	return s.nextSessionId
}

/**
* регистрация нового клиента.
 */
func (s *WsServer) Add(c *WsClientListener) {
	s.chanOp <- func(ws *WsServer) {
		if ws.working {
			c.SessionId = ws.findNextSessionId()
			ws.clientMap[c.SessionId] = c
			log.Printf("Added new client: %d (%d) [%s]\n", c.SessionId, len(ws.clientMap), c.RemoteAddr)
			if ws.NewClientHandler != nil {
				go ws.NewClientHandler(c)
			}
			//log.Println("Now", len(s.clients), "clients connected.")
		} else {
			c.Done()
		}
	}
}

func (s *WsServer) NotifyIfOnTopicSubscribed(topicId string, callback WsNewClientHandler) {
	s.chanOp <- func(ws *WsServer) {
		callbacks, _ := ws.subscribedClientHandler[topicId]
		ws.subscribedClientHandler[topicId] = append(callbacks, callback)
	}
}

//func (s *WsServer) Subscribe(c *Client, event string) {
//	s.chanOp <- func(ws *WsServer) {
//		if ws.working {
//			clients, _ := ws.subscribeCallbacks[event]
//			ws.subscribeCallbacks[event] = append(clients, c)
//
//			if ws.SubscribeClientHandler != nil {
//				go ws.SubscribeClientHandler(event, c)
//			}
//
//		}
//	}
//}

func (s *WsServer) Unsubscribe(c *WsClientListener, event string) {
	s.chanOp <- func(ws *WsServer) {
		clients, found := ws.subscribeCallbacks[event]
		if found {
			clients = deleteClient(clients, c)
			if len(clients) == 0 {
				delete(ws.subscribeCallbacks, event)
			} else {
				ws.subscribeCallbacks[event] = clients
			}
		}
	}

}

func (s *WsServer) Subscribe(c *WsClientListener, topicId string) {
	s.chanOp <- func(ws *WsServer) {
		if ws.working {
			clients, _ := ws.subscribeCallbacks[topicId]
			ws.subscribeCallbacks[topicId] = append(clients, c)

			if ws.SubscribeClientHandler != nil {
				go ws.SubscribeClientHandler(topicId, c)
			}

			callbacks, ok := ws.subscribedClientHandler[topicId]
			if ok {
				for _, callback := range callbacks {
					go callback(c)
				}
			}
		}
	}
}

func (s *WsServer) Del(c *WsClientListener) {
	s.chanOp <- func(ws *WsServer) {
		log.Printf("Remove client: %d (%d) [%s]\n", c.SessionId, len(ws.clientMap), c.RemoteAddr)
		delete(ws.clientMap, c.SessionId)
	}

}
func (s *WsServer) doneClients() {
	for op := range s.chanOp {
		op(s)
	}
	for _, element := range s.clientMap {
		element.Done()
	}

	s.clientMap = nil
	s.subscribedClientHandler = nil

}
func (s *WsServer) HasSubscribers(action string) bool {
	_, found := s.subscribeCallbacks[action]
	return found
}
func (s *WsServer) notifyClients(topicId string, code int16, data interface{}) {
	clients, found := s.subscribeCallbacks[topicId]
	if found {
		for _, client := range clients {
			client.Notify(topicId, &WsResponceMessage{
				Code: code,
				Data: data,
			})
		}
	}
}

func (s *WsServer) NotifyClients(topicId string, code int16, data interface{}) {
	s.chanOp <- func(ws *WsServer) {
		ws.notifyClients(topicId, code, data)
	}
}

func (s *WsServer) notifyClient(sessionId uint64, topicId string, code int16, data interface{}) {
	client, found := s.clientMap[sessionId]
	if found {
		client.Notify(topicId, &WsResponceMessage{
			Code: code,
			Data: data,
		})
	}
}

func (s *WsServer) NotifyClient(clientId uint64, topicId string, code int16, data interface{}) {
	if data != nil {
		s.chanOp <- func(s *WsServer) {
			s.notifyClient(clientId, topicId, code, data)
		}
	}
}

func (s *WsServer) Done() {
	s.working = false
	s.doneCh <- true
	s.doneClients()
}

func (s *WsServer) Err(err error) {
	s.errCh <- err
}
func (s *WsServer) Register(ws WsHandler) {
	wsCallbacks := ws.GetCallback()
	wsPrefix := ws.WsPrefix()

	for name, function := range wsCallbacks {
		s.callbacks[wsPrefix+"."+name] = function
		log.Println("Registered " + wsPrefix + "." + name)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *WsServer) Listen() {
	defer s.doneClients()

	for {
		select {
		// Add new a client
		case op := <-s.chanOp:
			op(s)

		case err := <-s.errCh:
			log.Println("WS Client Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

// websocket handler
func (s *WsServer) NewClientConnected(user WsUser, ws *websocket.Conn) {
	client := newWsClient(user, ws, s, ws.RemoteAddr())

	defer func() {
		err := ws.Close()
		if err != nil {
			s.errCh <- err
		}
	}()

	s.Add(client)
	client.Listen()
}
