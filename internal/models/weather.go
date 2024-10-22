package models

import (
	"database/sql"
	"time"
)

type Weather struct {
	SpotId           int       `db:"spot_id"`
	Time             time.Time `db:"timestamp"`
	AirTemperature   float64   `db:"air_temperature"`
	CurrentSpeed     float64   `db:"current_speed"`
	SeaLevel         float64   `db:"sea_level"`
	SwellDirection   float64   `db:"swell_direction"`
	SwellHeight      float64   `db:"swell_height"`
	SwellPeriod      float64   `db:"swell_period"`
	WaterTemperature float64   `db:"water_temperature"`
	WaveDirection    float64   `db:"wave_direction"`
	WaveHeight       float64   `db:"wave_height"`
	WavePeriod       float64   `db:"wave_period"`
	WindDirection    float64   `db:"wind_direction"`
	WindSpeed        float64   `db:"wind_speed"`
}

type WeatherModel struct {
	DB *sql.DB
}

func (w WeatherModel) GetWeatherDataFromDb(spotId int, start time.Time, duration int) ([]Weather, error) {
	rows, err := w.DB.Query(`
        SELECT spot_id, timestamp, air_temperature, current_speed, sea_level, swell_direction, swell_height, swell_period, water_temperature, wave_direction, wave_height, wave_period, wind_direction, wind_speed
        FROM weather
        WHERE spot_id = $1 AND timestamp BETWEEN $2 AND $3 AND EXTRACT(HOUR FROM timestamp) > 5 AND EXTRACT(HOUR FROM timestamp) <= 22
    `, spotId, start, start.Add(time.Duration(duration)*24*time.Hour))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weatherRows []Weather

	for rows.Next() {
		var weather Weather
		err := rows.Scan(
			&weather.SpotId,
			&weather.Time,
			&weather.AirTemperature,
			&weather.CurrentSpeed,
			&weather.SeaLevel,
			&weather.SwellDirection,
			&weather.SwellHeight,
			&weather.SwellPeriod,
			&weather.WaterTemperature,
			&weather.WaveDirection,
			&weather.WaveHeight,
			&weather.WavePeriod,
			&weather.WindDirection,
			&weather.WindSpeed,
		)
		if err != nil {
			return nil, err
		}
		weatherRows = append(weatherRows, weather)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return weatherRows, nil
}
