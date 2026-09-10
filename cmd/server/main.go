package main

import (
    "log"
    "net/http"

    "github.com/joho/godotenv" 
    "coding-profile-service/internal/cache"
    "coding-profile-service/internal/handler"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}