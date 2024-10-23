package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type SpotConfig struct {
	Id        int     `yaml:"id"`
	Name      string  `yaml:"name"`
	Lat       float64 `yaml:"latitude"`
	Long      float64 `yaml:"longitude"`
	Direction int     `yaml:"direction"`
}

type StormglassConfig struct {
	Url    string `yaml:"url"`
	ApiKey string `yaml:"api_key"`
}

type WeatherDataConfig struct {
	Source string `yaml:"source"`
}

type Config struct {
	Spots       []SpotConfig      `yaml:"spots"`
	Stormglass  StormglassConfig  `yaml:"stormglass"`
	WeatherData WeatherDataConfig `yaml:"weather_data"`
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return config
}

func SetConfig(cfg *Config) {
	once.Do(func() {
		config = cfg
	})
}
