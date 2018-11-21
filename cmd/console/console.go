package main

import (
	"fmt"
	sportsmonks "github.com/joesweeny/statshub/internal/service/sportmonks"
	"github.com/joesweeny/statshub/internal/config"
)

func main() {
	con := config.GetConfig()

	client := sportsmonks.NewClient(
		con.Services.SportMonks.BaseUri,
		con.Services.SportMonks.ApiKey,
	)

	response, err := client.GetCountries()

	if err != nil {
		panic(err.Error())
	}

	for _, country := range response.CountryList {
		fmt.Printf("%+v\n", country)
	}
}

