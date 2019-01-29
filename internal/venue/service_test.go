package venue

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
	t.Helper()
	venueRepo := new(mockVenueRepository)
	seasonRepo := new(mockSeasonRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/venues/season/123?api_token=my-key")
		b, _ := json.Marshal(venuesResponse())
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
		Repository: venueRepo,
		SeasonRepo: seasonRepo,
		Factory:    Factory{Clock: clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new venue", func(t *testing.T) {
		seasonRepo.On("Ids").Return([]int{123}, nil)
		venueRepo.On("GetById", 23).Return(&model.Venue{}, errors.New("not found"))
		venueRepo.On("Insert", mock.Anything).Return(nil)
		venueRepo.AssertNotCalled(t, "Update", mock.Anything)
		service.Process()
	})

	t.Run("updates existing venue", func(t *testing.T) {
		v := newVenue(34)
		seasonRepo.On("Ids").Return([]int{123}, nil)
		venueRepo.On("GetById", 23).Return(v, nil)
		venueRepo.On("Update", &v).Return(nil)
		venueRepo.AssertNotCalled(t, "Insert", mock.Anything)
		service.Process()
	})
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

type mockSeasonRepository struct {
	mock.Mock
}

func (m mockSeasonRepository) Insert(c *model.Season) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockSeasonRepository) Update(c *model.Season) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockSeasonRepository) Id(id int) (*model.Season, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Season)
	return c, args.Error(1)
}

func (m mockSeasonRepository) Ids() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m mockSeasonRepository) CurrentSeasonIds() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

type mockVenueRepository struct {
	mock.Mock
}

func (m mockVenueRepository) Insert(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m mockVenueRepository) Update(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m mockVenueRepository) GetById(id int) (*model.Venue, error) {
	args := m.Called(id)
	v := args.Get(0).(*model.Venue)
	return v, args.Error(1)
}

func venuesResponse() sportmonks.VenuesResponse {
	r := sportmonks.Venue{
		ID:       23,
		Name:     "London Stadium",
		Surface:  "Grass",
		City:     "London",
		Capacity: 60000,
		Image:    "/path/to/image",
	}

	res := sportmonks.VenuesResponse{}

	res.Data = append(res.Data, r)

	return res
}
