package round

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
	roundRepo := new(mockRoundRepository)
	seasonRepo := new(mockSeasonRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/seasons/100?api_token=my-key&include=rounds")
		b, _ := json.Marshal(seasonResponse())
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
		Repository: roundRepo,
		SeasonRepo: seasonRepo,
		Factory:    Factory{Clock: clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new round", func(t *testing.T) {
		done := make(chan bool)

		seasonRepo.On("Ids").Return([]int{100}, nil)
		roundRepo.On("GetById", 54).Return(&model.Round{}, errors.New("not found"))
		roundRepo.On("Insert", mock.Anything).Return(nil)
		roundRepo.AssertNotCalled(t, "Update", mock.Anything)
		service.Process("round", done)
	})

	t.Run("updates existing round", func(t *testing.T) {
		done := make(chan bool)

		r := newRound(34)
		seasonRepo.On("Ids").Return([]int{100}, nil)
		roundRepo.On("GetById", 34).Return(r, nil)
		roundRepo.On("Update", &r).Return(nil)
		roundRepo.AssertNotCalled(t, "Insert", mock.Anything)
		service.Process("round", done)
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

type mockRoundRepository struct {
	mock.Mock
}

func (m mockRoundRepository) Insert(c *model.Round) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRoundRepository) Update(c *model.Round) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockRoundRepository) GetById(id int) (*model.Round, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Round)
	return c, args.Error(1)
}

func seasonResponse() sportmonks.SeasonResponse {
	var round = 10
	var stage = 567
	s := sportmonks.Season{
		ID:              100,
		Name:            "2018-2019",
		LeagueID:        231,
		IsCurrentSeason: true,
		CurrentRoundID:  &round,
		CurrentStageID:  &stage,
		Fixtures: struct {
			Data []sportmonks.Fixture `json:"data"`
		}{},
		Rounds: struct {
			Data []sportmonks.Round `json:"data"`
		}{},
	}

	s.Rounds.Data = append(s.Rounds.Data, *newClientRound("2019-03-12", "2019-03-19"))

	res := sportmonks.SeasonResponse{}
	res.Data = s

	return res
}
