package country

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
	"net/http"
	"github.com/stretchr/testify/assert"
	"log"
	"io/ioutil"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/jonboulle/clockwork"
	"encoding/json"
	"errors"
	"bytes"
	"time"
)

func TestProcess(t *testing.T) {
	repo := new(mockRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
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
		repository: repo,
		factory:    factory{clockwork.NewFakeClockAt(date)},
		client:     client,
		logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new country", func (t *testing.T) {
		repo.On("getByExternalId", 1).Return(model.Country{}, errors.New("not Found"))
		repo.On("insert", mock.Anything).Return(nil)
		repo.AssertNotCalled(t, "update", mock.Anything)
		service.Process()
	})

	t.Run("updates existing country", func (t *testing.T) {
		id := uuid.FromStringOrNil("644e6546-63ae-47f8-be9d-06936e6fad35")

		repo.On("getByExternalId", 1).Return(newCountry(id), nil)

		repo.On("update", model.Country{}).Return(nil)
		repo.MethodCalled("update", model.Country{})
		repo.AssertNotCalled(t, "insert", mock.Anything)

		service.Process()
	})
}

type mockRepository struct {
	mock.Mock
}

func (m mockRepository) insert(c model.Country) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRepository) update(c model.Country) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRepository) getById(u uuid.UUID) (model.Country, error) {
	args := m.Called(u)
	return args.Get(0).(model.Country), args.Error(1)
}

func (m mockRepository) getByExternalId(id int) (model.Country, error) {
	args := m.Called(id)
	return args.Get(0).(model.Country), args.Error(1)
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
		ID:    1,
		Name:  "Brazil",
		Extra: struct {
			Continent string `json:"continent"`
			SubRegion string `json:"sub_region"`
			WorldRegion string `json:"world_region"`
			Fifa interface{} `json:"fifa,string"`
			ISO string `json:"iso"`
			Longitude string `json:"longitude"`
			Latitude string `json:"latitude"`
		} {
			Continent:   "South America",
			SubRegion:   "South America",
			WorldRegion: "South America",
			Fifa:        "BRA",
			ISO:         "BRA",
		},
	}
}

var date = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

