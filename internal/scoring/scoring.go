package scoring

import (
	"lr-surf-forecast/config"
	"lr-surf-forecast/internal/stormglass"
	"math"
)

// scale wage height to a value between 0 and 5
func scaleWaveHeight(waveHeight float64) float64 {
	idealWaveHeightMin := 0.8
	idealWaveHeightMax := 2.0

	if waveHeight < idealWaveHeightMin {
		return (waveHeight / idealWaveHeightMin) * 5
	} else if waveHeight > idealWaveHeightMax {
		score := 5 - ((waveHeight - idealWaveHeightMax) / idealWaveHeightMax * 5)
		if score < 0 {
			return 0
		}
		return score
	}
	return 5
}

// scale swell direction to a value between 0 and 5
func scaleSwellDirection(swellDirection float64, spotDirection int) float64 {
	directionDiff := math.Abs(swellDirection - float64(spotDirection))
	// Calculate direction difference (mod 360 to handle wrap-around cases)
	// 350° and 10° are only 20° apart for example
	if directionDiff > 180 {
		directionDiff = 360 - directionDiff
	}

	// 0° = 5, 90° = 0, 180° = -5
	directionScore := 5 - (directionDiff / 18)
	return math.Max(0, directionScore)
}

// scale swell period to a value between 0 and 5
func scaleSwellPeriod(swellPeriod float64) float64 {
	// Swell period scaling: Long periods (10s+) are usually better
	var periodScore float64
	if swellPeriod >= 10 {
		periodScore = 5
	} else {
		periodScore = swellPeriod / 2 // Scale period to 0-5 for periods less than 10s
	}
	return periodScore
}

func calculateSwellScore(swellHeight, swellPeriod, swellDirection float64, spot config.SpotConfig) float64 {
	directionScore := scaleSwellDirection(swellDirection, spot.Direction)
	periodScore := scaleSwellPeriod(swellPeriod)
	heightScore := scaleWaveHeight(swellHeight)

	swellScore := (0.4 * heightScore) + (0.4 * periodScore) + (0.2 * directionScore)
	return swellScore
}

// Function to calculate wind score based on speed and direction
func calculateWindScore(windSpeed, windDirection float64, spot config.SpotConfig) float64 {
	// Offshore wind is best, which occurs when wind blows from land to sea
	// Calculate angle difference between wind and coastline direction
	windDiff := math.Abs(windDirection - float64(spot.Direction))
	if windDiff > 180 {
		windDiff = 360 - windDiff
	}

	// 180° is perfect offshore
	windDirectionScore := 5 - (windDiff / 18)

	// Scale wind speed: Light winds (under 5 m/s) are ideal
	if windSpeed <= 5 {
		return windDirectionScore
	}
	return windDirectionScore - ((windSpeed - 5) / 2) // Penalty for high wind
}

// calculate comfort score based on water temperature and air temperature
func calculateComfort(waterTemperature, airTemperature float64) float64 {
	// we use 22 as the ideal temperature
	waterScore := 5 - math.Abs(22-waterTemperature)
	airScore := 5 - math.Abs(22-airTemperature)
	comfortScore := (0.5 * waterScore) + (0.5 * airScore)
	return comfortScore
}

func CalculateScoreSpotByHour(spot config.SpotConfig, hour stormglass.Hour) float64 {
	if hour.WaveHeight.Sg == 0.0 {
		return 0.0
	}
	waveScore := scaleWaveHeight(hour.WaveHeight.Sg)
	swellScore := calculateSwellScore(hour.SwellHeight.Sg, hour.SwellPeriod.Sg, hour.SwellDirection.Sg, spot)
	windScore := calculateWindScore(hour.WindSpeed.Sg, hour.WindDirection.Sg, spot)
	comfortScore := calculateComfort(hour.WaterTemperature.Sg, hour.AirTemperature.Sg)

	finalScore := (0.5 * waveScore) + (0.25 * swellScore) + (0.2 * windScore) + (0.05 * comfortScore)
	return finalScore
}
