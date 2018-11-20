package main

import (
	"os"
	"fmt"
	sportsmonks "github.com/joesweeny/statshub/console/service/sportmonks"
)

func main() {
	uri := os.Getenv("SPORT_MONKS_URI")
	key := os.Getenv("SPORT_MONKS_API_KEY")

	client := sportsmonks.NewClient(uri, key)

	response, err := client.GetCountries()

	if err != nil {
		panic(err.Error())
	}

	for _, country := range response.CountryList {
		fmt.Printf("%+v\n", country)
	}
}

