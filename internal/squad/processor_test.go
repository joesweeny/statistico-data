package squad

import (
	"bytes"
	"encoding/json"
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
	squadRepo := new(mockSquadRepository)
	seasonRepo := new(mockSeasonRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "http://example.com/api/v2.0/teams/season/100?api_token=my-key" {
			assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/teams/season/100?api_token=my-key")
			b, _ := json.Marshal(teamsResponse())
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			}
		}

		if req.URL.String() == "http://example.com/api/v2.0/squad/season/100/team/56?api_token=my-key" {
			assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/squad/season/100/team/56?api_token=my-key")
			b, _ := json.Marshal(squadResponse())
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			}
		}

		return &http.Response{}
	})

	client := sportmonks.Client{
		Client:  server,
		BaseURL: "http://example.com",
		ApiKey:  "my-key",
	}

	processor := Processor{
		Repository: squadRepo,
		SeasonRepo: seasonRepo,
		Factory:    Factory{Clock: clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new squad", func(t *testing.T) {
		done := make(chan bool)

		seasonRepo.On("Ids").Return([]int{100}, nil)
		squadRepo.On("BySeasonAndTeam", 100, 56).Return(&model.Squad{}, ErrNotFound)
		squadRepo.On("Insert", mock.Anything).Return(nil)
		squadRepo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("squad", done)
	})

	t.Run("updates existing squad", func(t *testing.T) {
		done := make(chan bool)

		squad := newSquad(34, 44)

		squadRepo.On("CurrentSeason").Return([]model.Squad{*squad}, nil)
		squadRepo.On("Update", mock.Anything).Return(nil)
		processor.Process("squad", done)
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

type mockSquadRepository struct {
	mock.Mock
}

func (m mockSquadRepository) Insert(c *model.Squad) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockSquadRepository) Update(c *model.Squad) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockSquadRepository) BySeasonAndTeam(seasonId, teamId int) (*model.Squad, error) {
	args := m.Called(seasonId, teamId)
	c := args.Get(0).(*model.Squad)
	return c, args.Error(1)
}

func (m mockSquadRepository) All() ([]model.Squad, error) {
	args := m.Called()
	return args.Get(0).([]model.Squad), args.Error(1)
}

func (m mockSquadRepository) CurrentSeason() ([]model.Squad, error) {
	args := m.Called()
	return args.Get(0).([]model.Squad), args.Error(1)
}

func squadResponse() sportmonks.SquadResponse {
	squad := newClientSquad(45, 43, 3)

	res := sportmonks.SquadResponse{}
	res.Data = *squad

	return res
}

func teamsResponse() sportmonks.TeamsResponse {
	t := &sportmonks.Team{
		ID:           56,
		LegacyID:     34,
		Name:         "West Ham United",
		ShortCode:    "WHU",
		CountryID:    8,
		NationalTeam: false,
		Founded:      1898,
		VenueID:      99,
	}

	res := sportmonks.TeamsResponse{}
	res.Data = append(res.Data, *t)

	return res
}
