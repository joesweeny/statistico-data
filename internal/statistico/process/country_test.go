package process

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/statistico/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestProcess(t *testing.T) {
	repo := new(mock.CountryRepository)

	server := mock.HttpClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/countries?api_token=my-key&page=1")
		b, _ := json.Marshal(countryResponse())
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
		}
	})

	client := sportmonks.Client{
		Client:  server,
		BaseURL: "http://example.com",
		ApiKey:  "my-key",
	}

	processor := Processor{
		Repository: repo,
		Factory:    Factory{clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new country", func(t *testing.T) {
		done := make(chan bool)

		repo.On("GetById", 1).Return(&model.Country{}, errors.New("not Found"))
		repo.On("Insert", mock.Anything).Return(nil)
		repo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("country", "",done)
	})

	t.Run("updates existing country", func(t *testing.T) {
		done := make(chan bool)

		c := newCountry(1)
		repo.On("GetById", 1).Return(c, nil)
		repo.On("Update", &c).Return(nil)
		repo.MethodCalled("Update", &c)
		repo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("country", "", done)
	})
}

func countryResponse() sportmonks.CountriesResponse {
	c := clientCountry()

	m := sportmonks.Meta{}
	m.Pagination.Total = 1
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := sportmonks.CountriesResponse{}
	res.Data = append(res.Data, c)
	res.Meta = m

	return res
}

func clientCountry() sportmonks.Country {
	return sportmonks.Country{
		ID:   1,
		Name: "Brazil",
		Extra: struct {
			Continent   string      `json:"continent"`
			SubRegion   string      `json:"sub_region"`
			WorldRegion string      `json:"world_region"`
			Fifa        interface{} `json:"fifa,string"`
			ISO         string      `json:"iso"`
			Longitude   string      `json:"longitude"`
			Latitude    string      `json:"latitude"`
		}{
			Continent:   "South America",
			SubRegion:   "South America",
			WorldRegion: "South America",
			Fifa:        "BRA",
			ISO:         "BRA",
		},
	}
}

