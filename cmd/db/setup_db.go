package main

import (
	"database/sql"
	"fmt"
	"go-surf-forecast/config"
	"go-surf-forecast/internal/stormglass"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func initSpotTable(db *sql.DB) {
	spotTable := `CREATE TABLE IF NOT EXISTS spot (
		spot_id SERIAL PRIMARY KEY,
		spot_name VARCHAR(255),
		latitude FLOAT,
		longitude FLOAT
	);`
	_, err := db.Exec(spotTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Spot table created successfully")
	}

	cfg := config.GetConfig()

	for _, spot := range cfg.Spots {
		_, err = db.Exec(`INSERT INTO spot (spot_id, spot_name, latitude, longitude) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (spot_id) DO NOTHING`,
			spot.Id, spot.Name, spot.Lat, spot.Long)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Spot data inserted successfully")
}

func initWeatherDataTable(db *sql.DB, dataSource string) {
	weatherTable := `CREATE TABLE IF NOT EXISTS weather (
        spot_id INT,
        timestamp TIMESTAMP,
        air_temperature FLOAT,
        current_speed FLOAT,
        sea_level FLOAT,
        swell_direction FLOAT,
        swell_height FLOAT,
        swell_period FLOAT,
        water_temperature FLOAT,
        wave_direction FLOAT,
        wave_height FLOAT,
        wave_period FLOAT,
        wind_direction FLOAT,
        wind_speed FLOAT,
		PRIMARY KEY (spot_id, timestamp),
		FOREIGN KEY (spot_id) REFERENCES spot(spot_id)
    );`
	_, err := db.Exec(weatherTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Weather table created successfully")
	}

	cfg := config.GetConfig()

	var weatherData *stormglass.StormglassWeatherPointApiResponse

	for _, spot := range cfg.Spots {

		duration := 7
		if dataSource == "stormglass" {
			start := time.Now()
			weatherData, err = stormglass.GetStormglassWeatherDataFromApi(spot, start, duration)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			start := time.Date(2024, time.October, 12, 0, 0, 0, 0, time.UTC)
			weatherData, err = stormglass.GetStormglassWeatherDataFromFile(spot, start, duration)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, data := range weatherData.Hours {
			_, err := db.Exec(`INSERT INTO weather(
            spot_id, timestamp, air_temperature, current_speed, sea_level, swell_direction, 
            swell_height, swell_period, water_temperature, wave_direction, wave_height, 
            wave_period, wind_direction, wind_speed) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (spot_id, timestamp) DO NOTHING`,
				spot.Id, data.Time, data.AirTemperature.Sg, data.CurrentSpeed.Sg, data.SeaLevel.Sg,
				data.SwellDirection.Sg, data.SwellHeight.Sg, data.SwellPeriod.Sg, data.WaterTemperature.Sg,
				data.WaveDirection.Sg, data.WaveHeight.Sg, data.WavePeriod.Sg, data.WindDirection.Sg, data.WindSpeed.Sg)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Weather data inserted successfully for spot", spot.Id)
	}
}

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	config.SetConfig(cfg)

	log.Println("Starting database setup...")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	weatherDataSource := os.Getenv("WEATHER_DATA_SOURCE")
	if weatherDataSource == "" {
		weatherDataSource = "file"
	}

	if postgresUser == "" || postgresPassword == "" || postgresDb == "" {
		log.Fatal("POSTGRES_USER, POSTGRES_PASSWORD and POSTGRES_DB must be set")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDb)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initSpotTable(db)
	log.Printf("Using data source = %s to init weather db...", weatherDataSource)
	initWeatherDataTable(db, weatherDataSource)
	log.Println("Database setup completed successfully.")

}
