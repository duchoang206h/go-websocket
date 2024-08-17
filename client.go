package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	ip := "198.19.249.2"
	exposePort := ":31066"
	connections := 10000
	u := url.URL{Scheme: "ws", Host: ip + exposePort, Path: "/ws"}
	log.Printf("Starting to connect to %s", u.String())
	log.Printf("Connecting to %s", u.String())
	var conns []*websocket.Conn
	for i := 0; i < connections; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Failed to connect", i, err)
			break
		}
		conns = append(conns, c)
		defer func() {
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			time.Sleep(time.Second)
			c.Close()
		}()
	}

	log.Printf("Finished initializing %d connections", len(conns))
	tts := time.Microsecond
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			log.Printf("Conn %d sending message", i)
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
				fmt.Printf("Failed to receive pong: %v", err)
			}
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from conn %v", i)))
		}
	}
}
