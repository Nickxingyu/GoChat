package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

var wsServer *WsServer

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func init() {
	wsServer = &WsServer{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func GetWsServer() *WsServer {
	return wsServer
}

func (s *WsServer) Run() {
	for {
		select {

		case client := <-s.register:
			s.handleRegister(client)

		case client := <-s.unregister:
			s.handleUnregister(client)

		case message := <-s.broadcast:
			s.handleBroadcast(message)

		}
	}
}

func (s *WsServer) ServeWs(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := NewClient(s, wsConn, name[0])

	s.register <- client
}

func (s *WsServer) handleRegister(client *Client) {
	s.clients[client.id] = client
	client.StartReadWriteLoop()
}

func (s *WsServer) handleBroadcast(message *Message) {

}

func (s *WsServer) handleUnregister(client *Client) {
	delete(s.clients, client.id)
}
