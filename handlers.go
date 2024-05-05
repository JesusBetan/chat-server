package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func createChat(w http.ResponseWriter, r *http.Request) {
	var chat Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	chat.Id = uuid.New().String()      // Genera un ID único para el chat
	chat.Clients = []*websocket.Conn{} // Inicializa la lista de clientes
	chats[chat.Id] = &chat

	json.NewEncoder(w).Encode(map[string]string{"id-chat": chat.Id})
}

func getChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]

	chat, ok := chats[chatID]
	if !ok {
		http.Error(w, "Chat no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func getChatIDs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id-user"]

	var chatIDs []string
	for id, chat := range chats {
		for _, msg := range chat.Messages {
			if msg.SenderId == userID {
				chatIDs = append(chatIDs, id)
				break
			}
		}
	}

	json.NewEncoder(w).Encode(chatIDs)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]

	chat, ok := chats[chatID]
	if !ok {
		http.Error(w, "Chat no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(chat.Messages)
}

func joinChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]
	var user struct {
		ID string `json:"id-user"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	chat, ok := chats[chatID]
	if !ok {
		http.Error(w, "Chat no encontrado", http.StatusNotFound)
		return
	}

	// Simula unirse al chat añadiendo un mensaje automático
	chat.Messages = append(chat.Messages, Message{
		Id:        uuid.New().String(),
		SenderId:  user.ID,
		Content:   "¡Me he unido al chat!",
		Timestamp: time.Now(),
	})

	notifyClients(chat.Clients, Message{
		Id:        uuid.New().String(),
		SenderId:  "system",
		Content:   fmt.Sprintf("Usuario %s se ha unido al chat.", user.ID),
		Timestamp: time.Now(),
	})
}
