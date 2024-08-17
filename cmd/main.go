package main

import (
	"log"
	"net/http"

	"websocket/config"
	"websocket/internal"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	redis := internal.NewRedis(internal.RedisOpt{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})
	hub := internal.NewHub(internal.HubOpt{
		Rdb:        redis,
		Broadcast:  make(chan []byte),
		Register:   make(chan *internal.Client),
		Unregister: make(chan *internal.Client),
		Clients:    map[*internal.Client]bool{},
	})
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWs(hub, w, r)
	})

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
