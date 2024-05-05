package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Chat struct {
	ID        string            `json:"id"`
	CreatorID string            `json:"creator_id"`
	Messages  []Message         `json:"messages"`
	Clients   []*websocket.Conn `json:"-"`
}

type Message struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var chats map[string]*Chat
