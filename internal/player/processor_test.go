package player

import (
	"bytes"
	"encoding/json"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/jonboulle/clockwork"
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
	squadRepo := new(mockSquadRepository)
	playerRepo := new(mockPlayerRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/players/5?api_token=my-key")
		b, _ := json.Marshal(playerResponse())
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
		Repository: playerRepo,
		SquadRepo:  squadRepo,
		Factory:    Factory{clockwork.NewFakeClock()},
		Client:     &client,
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	t.Run("inserts a new player if not already present", func(t *testing.T) {
		done := make(chan bool)

		squadRepo.On("All").Return(newSquad(), nil)
		playerRepo.On("Id", 5).Return(&model.Player{}, ErrNotFound)
		playerRepo.On("Insert", mock.Anything).Return(nil)
		processor.Process("player", done)
	})

	t.Run("player is not inserted if already present", func(t *testing.T) {
		done := make(chan bool)

		squadRepo.On("All").Return(newSquad(), nil)
		playerRepo.On("Id", 5).Return(newPlayer(5), nil)
		playerRepo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("player", done)
	})
}

type mockPlayerRepository struct {
	mock.Mock
}

func (m mockPlayerRepository) Insert(p *model.Player) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m mockPlayerRepository) Update(p *model.Player) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m mockPlayerRepository) Id(id int) (*model.Player, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Player), args.Error(1)
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

type roundTripFunc func(req *http.Request) *http.Response

func (r roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func playerResponse() sportmonks.PlayerResponse {
	p := newClientPlayer()
	r := sportmonks.PlayerResponse{}
	r.Data = *p
	return r
}

func newSquad() []model.Squad {
	var squads []model.Squad

	s := model.Squad{
		SeasonID:  45,
		TeamID:    98,
		PlayerIDs: []int{5},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	squads = append(squads, s)

	return squads
}
