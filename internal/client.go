package internal

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("message:", string(message))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("write error:", err)
			return
		}
	}
	err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		log.Println("write close message error:", err)
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), id: uuid.NewString()}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
