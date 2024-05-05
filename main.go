package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	chats = make(map[string]*Chat)
	router := mux.NewRouter()

	setupRoutes(router)

	fmt.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/chats", createChat).Methods("POST")
	router.HandleFunc("/chats/{id-chat}", getChat).Methods("GET")
	router.HandleFunc("/chats/{id-user}", getChatIDs).Methods("GET")
	router.HandleFunc("/chats/{id-chat}/messages", getMessages).Methods("GET")
	router.HandleFunc("/chats/{id-chat}/users", joinChat).Methods("POST")
	router.HandleFunc("/ws/{id-chat}", handleConnections)
}
