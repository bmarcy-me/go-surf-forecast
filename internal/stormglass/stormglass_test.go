package stormglass

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-surf-forecast/config"
	"go-surf-forecast/test"
)

func TestGetStormglassWeatherDataFromApi(t *testing.T) {

	cfg := config.Config{
		Stormglass: config.StormglassConfig{
			Url:    "http://fake-stormglass-url.com",
			ApiKey: "fake-api-key",
		},
	}

	spot := config.SpotConfig{
		Id:   1,
		Lat:  37.7749,
		Long: -122.4194,
	}

	start := time.Now()
	duration := 1

	mockResponse, err := test.TestData.ReadFile("data/mock-stormglass-api.json")
	if err != nil {
		t.Fatalf("Failed to read mock response file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override the URL in the config
	cfg.Stormglass.Url = server.URL
	config.SetConfig(&cfg)

	response, err := GetStormglassWeatherDataFromApi(spot, start, duration)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Hours[0].AirTemperature.Sg != 15.0 {
		t.Errorf("Expected air temperature 15.0, got %f", response.Hours[0].AirTemperature.Sg)
	}
}
