package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	chat, exists := chats[chatID]
	if !exists {
		http.Error(w, "Chat no encontrado", http.StatusNotFound)
		return
	}

	chat.Clients = append(chat.Clients, ws)
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		notifyClients(chat.Clients, msg)
	}
}

func notifyClients(clients []*websocket.Conn, msg Message) {
	for _, client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("error: %v", err)
			client.Close()
		}
	}
}
