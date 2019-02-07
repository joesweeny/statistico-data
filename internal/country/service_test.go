package country

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestProcess(t *testing.T) {
	repo := new(mockRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
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

	service := Service{
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
		service.Process("country", done)
	})

	t.Run("updates existing country", func(t *testing.T) {
		done := make(chan bool)

		c := newCountry(1)
		repo.On("GetById", 1).Return(c, nil)
		repo.On("Update", &c).Return(nil)
		repo.MethodCalled("Update", &c)
		repo.AssertNotCalled(t, "Insert", mock.Anything)
		service.Process("country", done)
	})
}

type mockRepository struct {
	mock.Mock
}

func (m mockRepository) Insert(c *model.Country) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRepository) Update(c *model.Country) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockRepository) GetById(id int) (*model.Country, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Country)
	return c, args.Error(1)
}

type roundTripFunc func(req *http.Request) *http.Response

func (r roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
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
