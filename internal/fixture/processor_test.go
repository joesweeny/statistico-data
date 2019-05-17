package fixture

import (
	"bytes"
	"encoding/json"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	t.Helper()
	fixtureRepo := new(mockFixtureRepository)
	seasonRepo := new(mockSeasonRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/seasons/123?api_token=my-key&include=fixtures")
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

	processor := Processor{
		Repository: fixtureRepo,
		SeasonRepo: seasonRepo,
		Factory:    Factory{Clock: clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts new fixture", func(t *testing.T) {
		done := make(chan bool)

		seasonRepo.On("Ids").Return([]int{123}, nil)
		fixtureRepo.On("ById", uint64(34)).Return(&model.Fixture{}, errors.New("not found"))
		fixtureRepo.On("Insert", mock.Anything).Return(nil)
		fixtureRepo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("fixture", "", done)
	})

	t.Run("updates existing fixture", func(t *testing.T) {
		done := make(chan bool)

		f := newFixture(34)
		seasonRepo.On("Ids").Return([]int{123}, nil)
		fixtureRepo.On("ById", uint64(34)).Return(f, nil)
		fixtureRepo.On("Update", &f).Return(nil)
		fixtureRepo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("fixture", "", done)
	})

	t.Run("inserts new fixture", func(t *testing.T) {
		done := make(chan bool)

		seasonRepo.On("CurrentSeasonIds").Return([]int{123}, nil)
		fixtureRepo.On("ById", uint64(34)).Return(&model.Fixture{}, errors.New("not found"))
		fixtureRepo.On("Insert", mock.Anything).Return(nil)
		fixtureRepo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("fixture", "", done)
	})

	t.Run("updates existing fixture", func(t *testing.T) {
		done := make(chan bool)

		f := newFixture(34)
		seasonRepo.On("CurrentSeasonIds").Return([]int{123}, nil)
		fixtureRepo.On("ById", 34).Return(f, nil)
		fixtureRepo.On("Update", &f).Return(nil)
		fixtureRepo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("fixture", "", done)
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

type mockFixtureRepository struct {
	mock.Mock
}

func (m mockFixtureRepository) Insert(c *model.Fixture) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockFixtureRepository) Update(c *model.Fixture) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockFixtureRepository) ById(id uint64) (*model.Fixture, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Fixture)
	return c, args.Error(1)
}

func (m mockFixtureRepository) Ids() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m mockFixtureRepository) IdsBetween(from, to time.Time) ([]int, error) {
	args := m.Called(from, to)
	return args.Get(0).([]int), args.Error(1)
}

func (m mockFixtureRepository) Between(from, to time.Time) ([]model.Fixture, error) {
	args := m.Called(from, to)
	return args.Get(0).([]model.Fixture), args.Error(1)
}

func (m mockFixtureRepository) ByTeamId(id int64, limit int32, before time.Time) ([]model.Fixture, error) {
	args := m.Called(id, limit, before)
	return args.Get(0).([]model.Fixture), args.Error(1)
}

func (m mockFixtureRepository) BySeasonId(id int64) ([]model.Fixture, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Fixture), args.Error(1)
}

func (m mockFixtureRepository) ByHomeAndAwayTeam(homeTeamId, awayTeamId uint64, limit uint32, before time.Time) ([]model.Fixture, error) {
	args := m.Called(homeTeamId, awayTeamId, limit, before)
	return args.Get(0).([]model.Fixture), args.Error(1)
}

func (m mockFixtureRepository) TeamIdsForSeason(seasonId uint64) ([]int, error) {
	args := m.Called(seasonId)
	return args.Get(0).([]int), args.Error(1)
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
	}

	s.Fixtures.Data = append(s.Fixtures.Data, *newClientFixture())

	res := sportmonks.SeasonResponse{}
	res.Data = s

	return res
}
