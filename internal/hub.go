package internal

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	rdb        *redis.Client
}
type HubOpt struct {
	Rdb        *redis.Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
}

func NewHub(opt HubOpt) *Hub {
	hub := &Hub{
		broadcast:  opt.Broadcast,
		register:   opt.Register,
		unregister: opt.Unregister,
		clients:    opt.Clients,
		rdb:        opt.Rdb,
	}
	go hub.subscribeToRedis()
	return hub
}

func (h *Hub) publishToRedis(message []byte) {
	err := h.rdb.Publish(context.Background(), "broadcast", message).Err()
	if err != nil {
		log.Println("Failed to publish message to Redis:", err)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("Client %s registered", client.id)
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.publishToRedis(message)
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) subscribeToRedis() {
	pubsub := h.rdb.Subscribe(context.Background(), "broadcast")
	ch := pubsub.Channel()

	for msg := range ch {
		h.mu.Lock()
		for client := range h.clients {
			select {
			case client.send <- []byte(msg.Payload):
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}
