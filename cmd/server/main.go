package main

import (
    "log"
    "net/http"

    "github.com/joho/godotenv" 
    "coding-profile-service/internal/cache"
    "coding-profile-service/internal/handler"
)

func main() {

    // Load .env file (only works locally, ignored if file missing)
    if err := godotenv.Load(); err != nil {
        log.Println("!  No .env file found — using system env vars (normal on Render)")
    }
    cache.Init()

    mux := http.NewServeMux()
    mux.HandleFunc("/stats", handler.StatsHandler)
    mux.HandleFunc("/", handler.RequestHandler)

    log.Println("⚙ Server running at http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}