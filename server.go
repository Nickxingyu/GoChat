package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	clients  map[string]*Client
	register chan *Client
}

var wsServer *WsServer

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func init() {
	wsServer = &WsServer{
		clients:  make(map[string]*Client),
		register: make(chan *Client),
	}
}

func GetWsServer() *WsServer {
	return wsServer
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

	client := NewClient(wsConn, name[0])

	s.register <- client
}
