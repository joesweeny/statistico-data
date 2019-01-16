package config

import "os"

type Config struct {
	Database
	Services
}

func GetConfig() (*Config) {
	config := Config{}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	config.Services.SportsMonks = SportsMonks{
		BaseUri: 	"https://soccer.sportmonks.com",
		ApiKey: 	os.Getenv("SPORTMONKS_API_KEY"),
	}

	return &config
}