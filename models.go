package main

import (
	"strings"
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
	Description         string            `json:"description"`
	LastMessageUsername string            `json:"lastMessageUsername"`
	LastMessage         string            `json:"lastMessage"`
	LastMessageDateTime time.Time         `json:"lastMessageDateTime"`
	Messages            []Message         `json:"messages"`
	Users               []*User           `json:"users"`
	Clients             []*websocket.Conn `json:"-"`
}

type Message struct {
	Id             string     `json:"id"`
	SenderId       string     `json:"senderId"`
	SenderUsername string     `json:"senderUsername"`
	Content        string     `json:"content"`
	DateTime       CustomTime `json:"dateTime"`
}

// CustomTime is a custom time.Time type that allows us to parse
// a custom time format
type CustomTime struct {
	time.Time
}

const ctLayout = "Jan 2, 2006 3:04:05 PM" // replace with your time format
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}
