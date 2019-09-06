package process

import (
	"bytes"
	"encoding/json"
	"errors"
	spClient "github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
	"github.com/stretchr/testify/assert"
	mk "github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestProcess(t *testing.T) {
	repo := new(mock.CountryRepository)
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

	requester := sportmonks.NewCountryRequester(&client)

	processor := CountryProcessor{
		repository: repo,
		requester:  requester,
		logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new country", func(t *testing.T) {
		done := make(chan bool)

		//a := assert.New(t)

		repo.On("GetById", 180).Return(&app.Country{}, errors.New("not Found"))
		repo.On("Insert", mk.AnythingOfType("app.Country")).Return(nil)
		repo.AssertCalled(t, "GetById", 180)
		repo.AssertCalled(t, "Insert", mk.AnythingOfType("app.Country"))
		repo.AssertNotCalled(t, "Update", mk.AnythingOfType("app.Country"))
		processor.Process("country", "", done)
	})

	//t.Run("updates existing country", func(t *testing.T) {
	//	done := make(chan bool)
	//
	//	c := stMock.Country(1)
	//	repo.On("GetById", 1).Return(c, nil)
	//	repo.On("Update", &c).Return(nil)
	//	repo.MethodCalled("Update", &c)
	//	repo.AssertNotCalled(t, "Insert", mock.Anything)
	//	processor.Process("country", "", done)
	//})
}

func countryResponse() *spClient.CountriesResponse {
	c := clientCountry()

	m := spClient.Meta{}
	m.Pagination.Total = 1
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := spClient.CountriesResponse{}
	res.Data = append(res.Data, *c)
	res.Meta = m

	return &res
}

func clientCountry() *spClient.Country {
	country := spClient.Country{
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


func newCountry(id int) *app.Country {
	c := app.Country{
		ID:        id,
		Name:      "England",
		Continent: "Europe",
		ISO:       "ENG",
	}

	return &c
}