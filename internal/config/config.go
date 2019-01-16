package config

import "os"

type Config struct {
	Database
	Services
}

func GetConfig() (*Config) {
	config := Config{}

	config.Database = Database{
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	}

	config.Services.SportsMonks = SportsMonks{
		baseUri: 	"https://soccer.sportmonks.com",
		apiKey: 	os.Getenv("SPORTMONKS_API_KEY"),
	}

	return &config
}