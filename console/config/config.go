package config

import "os"

type Config struct {
	Services Services
}

func GetConfig() (*Config) {
	sm := SportMonks{
		os.Getenv("SPORT_MONKS_URI"),
		os.Getenv("SPORT_MONKS_API_KEY"),
	}

	services := Services{sm}

	config := Config{services}

	return &config
}

type Services struct {
	SportMonks SportMonks
}

type SportMonks struct {
	BaseUri string
	ApiKey string
}
