package sportmonks

import (
	"bytes"
	"encoding/json"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestCountries(t *testing.T) {
	now := time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
	clock := clockwork.NewFakeClockAt(now)

	t.Run("parses response into Country struct and add pushes into channel", func (t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/countries?api_token=my-key&page=1")
			b, _ := json.Marshal(countryResponse())

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			}, nil
		})

		client := sportmonks.Client{
			Client:  server,
			BaseURL: "http://example.com",
			ApiKey:  "my-key",
		}

		service := CountryRequester{
			client: &client,
			clock: clock,
		}

		countries := make(chan *app.Country, 2)

		err := service.Countries(countries)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		for c := range countries {
			a.Equal(180, c.ID)
			a.Equal("England", c.Name)
			a.Equal("Europe", c.Continent)
			a.Equal("ENG", c.ISO)
		}
	})
}

func countryResponse() *sportmonks.CountriesResponse {
	c := clientCountry()

	m := sportmonks.Meta{}
	m.Pagination.Total = 1
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := sportmonks.CountriesResponse{}
	res.Data = append(res.Data, *c)
	res.Meta = m

	return &res
}

func clientCountry() *sportmonks.Country {
	country := sportmonks.Country{
		ID:   180,
		Name: "England",
		Extra: struct {
			Continent   string      `json:"continent"`
			SubRegion   string      `json:"sub_region"`
			WorldRegion string      `json:"world_region"`
			Fifa        interface{} `json:"fifa,string"`
			ISO         string      `json:"iso"`
			Longitude   string      `json:"longitude"`
			Latitude    string      `json:"latitude"`
		}{
			Continent:   "Europe",
			SubRegion:   "Western Europe",
			WorldRegion: "Europe",
			Fifa:        "ENG",
			ISO:         "ENG",
		},
	}

	return &country
}
