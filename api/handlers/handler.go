package handlers

import (
	"encoding/json"
	"fmt"
	"go-surf-forecast/config"
	"go-surf-forecast/internal/models"
	"go-surf-forecast/internal/scoring"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Spots []SurfSpot `json:"spots"`
}

type SurfSpot struct {
	Id      int              `json:"id"`
	Name    string           `json:"name"`
	Ratings []SurfSpotRating `json:"ratings"`
}

type SurfSpotRating struct {
	Rating float64   `json:"rating"`
	Time   time.Time `json:"time"`
}

var WeatherModel models.WeatherModel

func parseQueryParams(r *http.Request) (time.Time, int, error) {
	query := r.URL.Query()
	startParam := query.Get("start")
	durationParam := query.Get("duration")

	var start time.Time
	var err error
	if startParam == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startParam)
		if err != nil {
			return time.Time{}, 0, err
		}
	}

	duration := 7
	if durationParam != "" {
		duration, err = strconv.Atoi(durationParam)
		if duration < 1 || duration > 7 {
			return time.Time{}, 0, fmt.Errorf("duration must be between 1 and 7")
		}
		if err != nil {
			return time.Time{}, 0, err
		}
	}

	return start, duration, nil
}

// map weather data from database to API response
func weatherDataToApi(spotConfig config.SpotConfig, weatherData []models.Weather) SurfSpot {
	spot := SurfSpot{
		Id:   spotConfig.Id,
		Name: spotConfig.Name,
	}
	for _, weather := range weatherData {
		rating := SurfSpotRating{
			Rating: scoring.CalculateScoreSpotByHour(spotConfig, weather),
			Time:   weather.Time,
		}
		spot.Ratings = append(spot.Ratings, rating)
	}

	return spot
}

func getBestSpotAtAnytime(spots []SurfSpot) SurfSpot {
	var bestSpot SurfSpot
	var highestScore float64

	for _, spot := range spots {
		for _, rating := range spot.Ratings {
			if rating.Rating > highestScore {
				highestScore = rating.Rating
				bestSpot = spot
				bestSpot.Ratings = []SurfSpotRating{rating} // Keep only the best rating
			}
		}
	}

	return bestSpot
}

func Healtcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetSpots is a handler function that returns score for all spots by hour
func GetSpots(w http.ResponseWriter, r *http.Request) {
	start, duration, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	cfg := config.GetConfig()
	var response Response
	for _, spot := range cfg.Spots {
		weatherData, err := WeatherModel.GetWeatherDataFromDb(spot.Id, start, duration)
		if err != nil {
			http.Error(w, "Could not get static data", http.StatusInternalServerError)
			return
		}

		spotData := weatherDataToApi(spot, weatherData)
		response.Spots = append(response.Spots, spotData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBestSpot is a handler function that returns the spot with the best score at any time
func GetBestSpot(w http.ResponseWriter, r *http.Request) {
	start, duration, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var spots []SurfSpot
	cfg := config.GetConfig()
	for _, spotConfig := range cfg.Spots {
		weatherData, err := WeatherModel.GetWeatherDataFromDb(spotConfig.Id, start, duration)
		if err != nil {
			http.Error(w, "Could not get static data", http.StatusInternalServerError)
			return
		}

		spot := weatherDataToApi(spotConfig, weatherData)
		spots = append(spots, spot)
	}

	bestSpot := getBestSpotAtAnytime(spots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestSpot)
}
