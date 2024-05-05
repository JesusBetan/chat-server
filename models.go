package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	Id                  string            `json:"id"`
	ShortId             string            `json:"short_id"`
	CreatorId           string            `json:"creator_id"`
	Name                string            `json:"name"`
	LastMessageUsername string            `json:"last_message_username"`
	LastMessage         string            `json:"last_message"`
	LastMessageDateTime time.Time         `json:"last_message_datetime"`
	Messages            []Message         `json:"messages"`
	Users               []*User           `json:"users"`
	Clients             []*websocket.Conn `json:"-"`
}

type Message struct {
	Id        string    `json:"id"`
	SenderId  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var users map[string]*User

// make mock chats for testing
// var chats map[string]*Chat
var chats = map[string]*Chat{
	"1": {
		Id:                  "1",
		ShortId:             "ABC1",
		CreatorId:           "1",
		Name:                "Chat 1",
		LastMessageUsername: "User 1",
		LastMessage:         "Hola",
		LastMessageDateTime: time.Now(),
		Messages:            []Message{},
		Users:               []*User{{Id: "1", Username: "User 1"}},
		Clients:             []*websocket.Conn{},
	},
	"2": {
		Id:                  "2",
		ShortId:             "ABC2",
		CreatorId:           "2",
		Name:                "Chat 2",
		LastMessageUsername: "User 2",
		LastMessage:         "Hola",
		LastMessageDateTime: time.Now(),
		Messages:            []Message{},
		Users:               []*User{{Id: "1", Username: "User 1"}},
		Clients:             []*websocket.Conn{},
	},
}
