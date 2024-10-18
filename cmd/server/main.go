package main

import (
	"log"
	"net/http"

	"lr-surf-forecast/api/handlers"
)

func main() {
	http.HandleFunc("/api/spots", handlers.GetSpots)
	http.HandleFunc("/api/spots/best", handlers.GetBestSpot)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
