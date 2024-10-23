package stormglass

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"go-surf-forecast/config"
)

type StormglassWeatherPointApiResponse struct {
	Hours []Hour `json:"hours"`
	Meta  Meta   `json:"meta"`
}

type Hour struct {
	AirTemperature   Source    `json:"airTemperature"`
	CurrentSpeed     Source    `json:"currentSpeed"`
	SeaLevel         Source    `json:"seaLevel"`
	SwellDirection   Source    `json:"swellDirection"`
	SwellHeight      Source    `json:"swellHeight"`
	SwellPeriod      Source    `json:"swellPeriod"`
	Time             time.Time `json:"time"`
	WaterTemperature Source    `json:"waterTemperature"`
	WaveDirection    Source    `json:"waveDirection"`
	WaveHeight       Source    `json:"waveHeight"`
	WavePeriod       Source    `json:"wavePeriod"`
	WindDirection    Source    `json:"windDirection"`
	WindSpeed        Source    `json:"windSpeed"`
}

type Source struct {
	Sg float64 `json:"sg"`
}

type Meta struct {
	Cost         int      `json:"cost,omitempty"`
	DailyQuota   int      `json:"dailyQuota,omitempty"`
	End          string   `json:"end,omitempty"`
	Lat          float64  `json:"lat,omitempty"`
	Lng          float64  `json:"lng,omitempty"`
	Params       []string `json:"params,omitempty"`
	RequestCount int      `json:"requestCount,omitempty"`
	Start        string   `json:"start,omitempty"`
}

// call stormglass api endpoint v2/weather/point
func GetStormglassWeatherDataFromApi(spot config.SpotConfig, start time.Time, duration int) (*StormglassWeatherPointApiResponse, error) {
	cfg := config.GetConfig()
	stormglassApiKey := cfg.Stormglass.ApiKey
	if stormglassApiKey == "" {
		return nil, fmt.Errorf("stormglass API key is not set in config")
	}

	baseURL, err := url.Parse(cfg.Stormglass.Url)
	if err != nil {
		return nil, err
	}
	baseURL.Path += "/weather/point"

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", spot.Lat))
	params.Add("lng", fmt.Sprintf("%f", spot.Long))
	params.Add("params", "airTemperature,currentSpeed,seaLevel,swellDirection,swellHeight,swellPeriod,waterTemperature,waveDirection,waveHeight,wavePeriod,windDirection,windSpeed")
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	end := start.Add(time.Duration(duration) * 24 * time.Hour).Unix()
	params.Add("end", fmt.Sprintf("%d", end))
	params.Add("source", "sg")
	baseURL.RawQuery = params.Encode()

	log.Default().Printf("Calling stormglass API for spot %d", spot.Id)

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", stormglassApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherPointApiResponse StormglassWeatherPointApiResponse
	if err := json.Unmarshal(body, &weatherPointApiResponse); err != nil {
		return nil, err
	}

	return &weatherPointApiResponse, nil
}

// reads a static JSON file for a spot and returns the data
func GetStormglassWeatherDataFromFile(spot config.SpotConfig, start time.Time, duration int) (*StormglassWeatherPointApiResponse, error) {
	filePath := fmt.Sprintf("assets/data/stormglass-data-spot-%d.json", spot.Id)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return &StormglassWeatherPointApiResponse{}, err
	}

	var stormglassResponse StormglassWeatherPointApiResponse
	err = json.Unmarshal(file, &stormglassResponse)
	if err != nil {
		return &StormglassWeatherPointApiResponse{}, err
	}

	// Filter the Hours field
	endTime := start.Add(time.Duration(duration) * 24 * time.Hour)
	var filteredHours []Hour
	for _, hour := range stormglassResponse.Hours {
		hourTime := hour.Time
		if hourTime.After(start) && hourTime.Before(endTime) {
			// Keep only daylights hours
			if hourTime.Hour() > 5 && hourTime.Hour() <= 22 {
				filteredHours = append(filteredHours, hour)
			}
		}
	}

	stormglassResponse.Hours = filteredHours
	return &stormglassResponse, nil
}
