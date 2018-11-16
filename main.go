package main

import (
	"net/http"
	"fmt"
	//"log"
	"io/ioutil"
	"encoding/json"
)

type SportMonkCountry struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Extra struct {
		Continent string `json:"continent"`
		IsoCode string `json:"iso"`
	}
}

type SportMonkCountryResponse struct {
	Countries []SportMonkCountry `json:"data"`
}

func main() {
	response, err := http.Get("https://soccer.sportmonks.com/api/v2.0/countries?api_token=hMNoq0c2fMjipNWEeG7IMmDF9bMNKeVoRi8lJ0qZDhg125U1IormejZKfwua")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	r, err := parseCountries(data)

	for _, country := range r.Countries {
		fmt.Printf("%+v\n", country)
	}
}

func parseCountries(body []byte) (*SportMonkCountryResponse, error) {
	var r = new(SportMonkCountryResponse)

	err := json.Unmarshal(body, r)

	if err != nil {
		fmt.Printf("Parsing request body in objects fields: %s\n", err)
	}

	return r, err
}