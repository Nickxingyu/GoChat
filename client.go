package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	readWait = 60 * time.Second

	pingPeriod = (readWait * 9) / 10

	maxMessageSize = 10000
)

type Client struct {
	id       string
	conn     *websocket.Conn
	name     string
	wsServer *WsServer
	send     chan *Message
}

func NewClient(wsServer *WsServer, conn *websocket.Conn, name string) *Client {
	client := Client{
		id:       uuid.NewString(),
		conn:     conn,
		name:     name,
		wsServer: wsServer,
		send:     make(chan *Message),
	}
	return &client
}

func (c *Client) StartReadWriteLoop() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) readLoop() {
	defer func() {
		c.disconnect()
	}()

	c.resetReadDeadLine()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.resetReadDeadLine()
		return nil
	})

	for {
		var message = Message{}
		if err := c.conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("Unexpected Close Error: " + err.Error())
				break
			}

			c.handleMessage(&message)
		}
	}
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.resetWriteDeadLine()
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			} else {
				c.conn.WriteJSON(message)
			}

		case <-ticker.C:
			c.resetWriteDeadLine()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (c *Client) disconnect() {
	c.wsServer.unregister <- c
	close(c.send)
	c.conn.Close()
}

func (c *Client) resetReadDeadLine() {
	c.conn.SetReadDeadline(time.Now().Add(readWait))
}

func (c *Client) resetWriteDeadLine() {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
}

func (c *Client) handleMessage(message *Message) {
	message.SenderID = c.id
	c.wsServer.broadcast <- message
}
