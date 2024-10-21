package main

import (
	"log"
	"net/http"

	"go-surf-forecast/api/handlers"
	"go-surf-forecast/config"
)

func main() {

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	config.SetConfig(cfg)

	http.HandleFunc("/api/spots", handlers.GetSpots)
	http.HandleFunc("/api/spots/best", handlers.GetBestSpot)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
