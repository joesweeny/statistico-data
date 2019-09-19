package sportmonks_test

import (
	"bytes"
	"encoding/json"
	spClient "github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestCountries(t *testing.T) {
	t.Run("countries returns channel", func (t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/countries?api_token=my-key&page=1")
			b, _ := json.Marshal(countryResponse())

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			}, nil
		})

		client := spClient.Client{
			Client:  server,
			BaseURL: "http://example.com",
			ApiKey:  "my-key",
		}

		requester := sportmonks.NewCountryRequester(&client, &log.Logger{})

		ch := requester.Countries()

		eng := <- ch
		ger := <- ch

		a := assert.New(t)

		a.Equal(180, eng.ID)
		a.Equal("England", eng.Name)
		a.Equal("Europe", eng.Continent)
		a.Equal("ENG", eng.ISO)
		a.Equal(5, ger.ID)
		a.Equal("Germany", ger.Name)
		a.Equal("Europe", ger.Continent)
		a.Equal("ENG", ger.ISO)
	})
}

func countryResponse() *spClient.CountriesResponse {
	eng := clientCountry(180, "England")
	ger := clientCountry(5, "Germany")

	m := spClient.Meta{}
	m.Pagination.Total = 2
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := spClient.CountriesResponse{}
	res.Data = append(res.Data, *eng, *ger)
	res.Meta = m

	return &res
}

func clientCountry(id int, name string) *spClient.Country {
	country := spClient.Country{
		ID:   id,
		Name: name,
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
