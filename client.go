package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	conn *websocket.Conn
	name string
}

func NewClient(conn *websocket.Conn, name string) *Client {
	client := Client{
		id:   uuid.NewString(),
		conn: conn,
		name: name,
	}
	return &client
}

func (c *Client) StartReadWriteLoop() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) readLoop() {

}

func (c *Client) writeLoop() {

}