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

var users map[string]*User

var chats map[string]*Chat

func getChatsMock() map[string]*Chat {
	return map[string]*Chat{
		"ABC1": {
			Id:                  "1",
			ShortId:             "ABC1",
			CreatorId:           "1",
			Name:                "Chat 1",
			Description:         "Descripción del chat 1",
			LastMessageUsername: "User 1",
			LastMessage:         "Hola",
			LastMessageDateTime: time.Now(),
			Messages:            []Message{},
			Users:               []*User{{Id: "123456", Username: "User 1"}},
			Clients:             []*websocket.Conn{},
		},
		"ABC2": {
			Id:                  "2",
			ShortId:             "ABC2",
			CreatorId:           "2",
			Name:                "Chat 2",
			Description:         "Descripción del chat 2",
			LastMessageUsername: "User 2",
			LastMessage:         "Hola",
			LastMessageDateTime: time.Now(),
			Messages:            []Message{},
			Users:               []*User{{Id: "1", Username: "User 1"}},
			Clients:             []*websocket.Conn{},
		},
	}
}

func createChat(w http.ResponseWriter, r *http.Request) {
	var chat Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	//get user by id
	user := getUserById(chat.CreatorId)

	var chatUsers = []*User{user}

	chat.Id = uuid.New().String()      // Genera un ID único para el chat
	chat.ShortId = chat.Id[:6]         // Genera un ID corto para el chat
	chat.Messages = []Message{}        // Inicializa la lista de mensajes
	chat.Users = chatUsers             // Inicializa la lista de usuarios
	chat.Clients = []*websocket.Conn{} // Inicializa la lista de clientes

	chats[chat.ShortId] = &chat

	json.NewEncoder(w).Encode(chat)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	user.Id = uuid.New().String()
	users[user.Id] = &user

	json.NewEncoder(w).Encode(user)
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

	var userChats []*Chat = []*Chat{}
	for _, chat := range chats {
		for _, user := range chat.Users {
			if user.Id == userID {
				userChats = append(userChats, chat)
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

func createMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	message.Id = uuid.New().String()

	chat, ok := chats[chatID]
	if !ok {
		http.Error(w, "Chat no encontrado", http.StatusNotFound)
		return
	}

	chat.LastMessageUsername = getUserById(message.SenderId).Username
	chat.LastMessage = message.Content
	chat.Messages = append(chat.Messages, message)

	notifyClients(chat.Clients, message)
}

func joinChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id-chat"]
	var user User
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
	chat.Users = append(chat.Users, getUserById(user.Id))

	// Añadir un mensaje automático
	chat.Messages = append(chat.Messages, Message{
		Id:             uuid.New().String(),
		SenderId:       user.Id,
		SenderUsername: user.Username,
		Content:        "¡Me he unido al chat!",
		DateTime:       CustomTime{time.Now()},
	})

	notifyClients(chat.Clients, Message{
		Id:       uuid.New().String(),
		SenderId: "system",
		Content:  fmt.Sprintf("Usuario %s se ha unido al chat.", user.Id),
		DateTime: CustomTime{time.Now()},
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
