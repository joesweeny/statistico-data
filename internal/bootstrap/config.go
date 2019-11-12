package bootstrap

import (
	"os"
)

type Config struct {
	Database
	Services
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Services struct {
	SportsMonks
}

type SportsMonks struct {
	ApiKey  string
}


func GetConfig() *Config {
	config := Config{}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	config.SportsMonks = SportsMonks{
		ApiKey:  os.Getenv("SPORTMONKS_API_KEY"),
	}

	return &config
}
