package team

import (
	"bytes"
	"encoding/json"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestProcess(t *testing.T) {
	t.Helper()
	teamRepo := new(mockTeamRepository)
	seasonRepo := new(mockSeasonRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/teams/season/100?api_token=my-key")
		b, _ := json.Marshal(teamsResponse())
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
		Repository: teamRepo,
		SeasonRepo: seasonRepo,
		Factory:    Factory{Clock: clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new round", func(t *testing.T) {
		done := make(chan bool)

		seasonRepo.On("Ids").Return([]int{100}, nil)
		teamRepo.On("GetById", 56).Return(&model.Team{}, ErrNotFound)
		teamRepo.On("Insert", mock.Anything).Return(nil)
		teamRepo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("team", "", done)
	})

	t.Run("updates existing round", func(t *testing.T) {
		done := make(chan bool)

		r := newTeam(34)
		seasonRepo.On("Ids").Return([]int{100}, nil)
		teamRepo.On("GetById", 34).Return(r, nil)
		teamRepo.On("Update", &r).Return(nil)
		teamRepo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("team", "", done)
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

type mockTeamRepository struct {
	mock.Mock
}

func (m mockTeamRepository) Insert(c *model.Team) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockTeamRepository) Update(c *model.Team) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockTeamRepository) GetById(id int) (*model.Team, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Team)
	return c, args.Error(1)
}

func teamsResponse() sportmonks.TeamsResponse {
	team := newClientTeam()

	res := sportmonks.TeamsResponse{}
	res.Data = append(res.Data, *team)

	return res
}
