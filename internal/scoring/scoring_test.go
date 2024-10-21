package scoring

import (
	"go-surf-forecast/config"
	"go-surf-forecast/internal/stormglass"
	"testing"
)

func TestScaleWaveHeight(t *testing.T) {
	testCases := []struct {
		waveHeight float64
		expected   float64
	}{
		{0.0, 0.0},
		{0.2, 1.25},
		{0.4, 2.5},
		{0.8, 5.0},
		{1.4, 5.0},
		{2.0, 5.0},
		{2.1, 4.75},
		{10.0, 0.0},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Logf("Testing wave height scale %f", tc.waveHeight)
			result := scaleWaveHeight(tc.waveHeight)
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestScaleSwellDirection(t *testing.T) {
	testCases := []struct {
		swellDirection float64
		spotDirection  int
		expected       float64
	}{
		{0.0, 0, 5.0},
		{90.0, 0, 0.0},
		{180.0, 180.0, 5.0},
		{270.0, 0, 0.0},
		{350.0, 10.0, 3.888888888888889},
		{350.0, 20.0, 3.333333333333333},
		{10.0, 350.0, 3.888888888888889},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Logf("Testing swell direction scale %f against spot %d", tc.swellDirection, tc.spotDirection)
			result := scaleSwellDirection(tc.swellDirection, tc.spotDirection)
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestCalculateScoreSpotByHour(t *testing.T) {
	testCases := []struct {
		spot     config.SpotConfig
		hour     stormglass.Hour
		label    string
		expected float64
	}{
		{
			spot: config.SpotConfig{Direction: 90},
			hour: stormglass.Hour{
				WaveHeight:       stormglass.Source{Sg: 1.0},
				SwellHeight:      stormglass.Source{Sg: 1.0},
				SwellPeriod:      stormglass.Source{Sg: 10.0},
				SwellDirection:   stormglass.Source{Sg: 90.0},
				WindSpeed:        stormglass.Source{Sg: 4.0},
				WindDirection:    stormglass.Source{Sg: 90.0},
				WaterTemperature: stormglass.Source{Sg: 22.0},
				AirTemperature:   stormglass.Source{Sg: 22.0},
			},
			label:    "perfect conditions",
			expected: 5.0,
		},
		{
			spot: config.SpotConfig{Direction: 0},
			hour: stormglass.Hour{
				WaveHeight:       stormglass.Source{Sg: 0.0},
				SwellHeight:      stormglass.Source{Sg: 2.0},
				SwellPeriod:      stormglass.Source{Sg: 8.0},
				SwellDirection:   stormglass.Source{Sg: 90.0},
				WindSpeed:        stormglass.Source{Sg: 15.0},
				WindDirection:    stormglass.Source{Sg: 270.0},
				WaterTemperature: stormglass.Source{Sg: 20.0},
				AirTemperature:   stormglass.Source{Sg: 18.0},
			},
			label:    "no wave",
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Logf("Testing global score by hour for %s", tc.label)
			result := CalculateScoreSpotByHour(tc.spot, tc.hour)
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}
