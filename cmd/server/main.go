package main

import (
	"log"
	"net/http"
	"coding-profile-service/internal/handler"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", handler.StatsHandler)

	mux.HandleFunc("/", handler.RequestHandler)

	log.Println("⚙ Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
