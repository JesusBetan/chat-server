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

	//get user by id
	user := getUserById(chat.CreatorId)

	chat.Id = uuid.New().String()      // Genera un ID único para el chat
	chat.ShortId = chat.Id[:4]         // Genera un ID corto para el chat
	chat.Messages = []Message{}        // Inicializa la lista de mensajes
	chat.Users = []*User{user}         // Inicializa la lista de usuarios
	chat.Clients = []*websocket.Conn{} // Inicializa la lista de clientes

	chats[chat.ShortId] = &chat

	json.NewEncoder(w).Encode(map[string]string{"id-chat": chat.ShortId})
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

func getUserChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id-user"]

	var userChats []string
	for _, chat := range chats {
		for _, user := range chat.Users {
			if user.Id == userID {
				userChats = append(userChats, chat.ShortId)
				break
			}
		}
	}

	json.NewEncoder(w).Encode(userChats)
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

	// Unir al usuario al chat
	chat.Users = append(chat.Users, getUserById(user.ID))

	// Añadir un mensaje automático
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

func getUserById(id string) *User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}
