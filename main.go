package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Call the next handler in the chain
		log.Printf("%s %s %s - %s\n", r.Method, r.URL.Path, r.Proto, time.Since(start))
	})
}

func main() {
	chats = getChatsMock()
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	setupRoutes(router)

	fmt.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/{id-user}/chats", getUserChats).Methods("GET")
	router.HandleFunc("/chats", createChat).Methods("POST")
	router.HandleFunc("/chats/{id-chat}", getChat).Methods("GET")
	router.HandleFunc("/chats/{id-chat}/messages", getMessages).Methods("GET")
	router.HandleFunc("/chats/{id-chat}/users", joinChat).Methods("POST")
	router.HandleFunc("/ws/{id-chat}", handleConnections)
}
