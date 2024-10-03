package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var streamHub *StreamHub

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new router
	r := mux.NewRouter()

	// Initialize streamHub
	streamHub = newStreamHub()
	go streamHub.run()

	r.HandleFunc("/webhook", handleWebhook).Methods("POST")

	// Add WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(streamHub, w, r)
	})

	// Set up CORS
	r.Use(mux.CORSMethodMiddleware(r))

	// Set up JSON parsing middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3005" // Default port if not specified in .env
	}
	fmt.Printf("Server is listening at Your Server Endpoint:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
