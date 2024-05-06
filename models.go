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
	ShortId             string            `json:"shortId"`
	CreatorId           string            `json:"creatorId"`
	Name                string            `json:"name"`
	LastMessageUsername string            `json:"lastMessageUsername"`
	LastMessage         string            `json:"lastMessage"`
	LastMessageDateTime time.Time         `json:"lastMessageDateTime"`
	Messages            []Message         `json:"messages"`
	Users               []*User           `json:"users"`
	Clients             []*websocket.Conn `json:"-"`
}

type Message struct {
	Id        string    `json:"id"`
	SenderId  string    `json:"senderId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
