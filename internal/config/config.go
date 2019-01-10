package config

import "os"

type Config struct {
	DB Database
}

func GetConfig() (*Config) {
	config := Config{}

	config.DB = Database{
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	}

	return &config
}

type Database struct {
	Driver    string
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
}
