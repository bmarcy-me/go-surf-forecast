package handlers

import (
	"encoding/json"
	"fmt"
	"lr-surf-forecast/config"
	"lr-surf-forecast/internal/scoring"
	"lr-surf-forecast/internal/stormglass"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Spots []SurfSpot `json:"spots"`
}

type SurfSpot struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	Ratings []SurfSpotRating `json:"ratings"`
}

type SurfSpotRating struct {
	Rating float64   `json:"rating"`
	Time   time.Time `json:"time"`
}

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

// map stormglass response to API response
func stormglassToApi(spotId int, stormglassResponse stormglass.StormglassWeatherPointApiResponse) SurfSpot {
	spot := SurfSpot{
		Id:   strconv.Itoa(spotId),
		Name: config.SpotConfigs[spotId-1].Name,
	}
	for _, hour := range stormglassResponse.Hours {
		rating := SurfSpotRating{
			Rating: scoring.CalculateScoreSpotByHour(config.SpotConfigs[spotId-1], hour),
			Time:   hour.Time,
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

// GetSpots is a handler function that returns score for all spots by hour
func GetSpots(w http.ResponseWriter, r *http.Request) {
	start, duration, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var response Response
	for _, spot := range config.SpotConfigs {
		stormglassApiResponse, err := stormglass.GetStaticStormglassData(spot.Id, start, duration)
		if err != nil {
			http.Error(w, "Could not get static data", http.StatusInternalServerError)
			return
		}

		spotData := stormglassToApi(spot.Id, *stormglassApiResponse)
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
	for _, spotConfig := range config.SpotConfigs {
		stormglassApiResponse, err := stormglass.GetStaticStormglassData(spotConfig.Id, start, duration)
		if err != nil {
			http.Error(w, "Could not get static data", http.StatusInternalServerError)
			return
		}

		spot := stormglassToApi(spotConfig.Id, *stormglassApiResponse)
		spots = append(spots, spot)
	}

	bestSpot := getBestSpotAtAnytime(spots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestSpot)
}
