package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-surf-forecast/api/handlers"
	"go-surf-forecast/config"
	"go-surf-forecast/internal/models"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	config.SetConfig(cfg)

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDb)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	handlers.WeatherModel = models.WeatherModel{DB: db}

	http.HandleFunc("/api/healthcheck", handlers.Healtcheck)
	http.HandleFunc("/api/spots", handlers.GetSpots)
	http.HandleFunc("/api/spots/best", handlers.GetBestSpot)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
