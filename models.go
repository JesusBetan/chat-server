package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Chat struct {
	Id                  string            `json:"id"`
	CreatorId           string            `json:"creator_id"`
	Name                string            `json:"name"`
	LastMessageUsername string            `json:"last_message_username"`
	LastMessage         string            `json:"last_message"`
	LastMessageDateTime time.Time         `json:"last_message_datetime"`
	Messages            []Message         `json:"messages"`
	Clients             []*websocket.Conn `json:"-"`
}

type Message struct {
	Id        string    `json:"id"`
	SenderId  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var chats map[string]*Chat
