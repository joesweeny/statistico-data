package competition

import (
	"bytes"
	"encoding/json"
	"errors"
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
	repo := new(mockRepository)

	server := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/leagues?api_token=my-key&page=1")
		b, _ := json.Marshal(leagueResponse())
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

	t.Run("inserts new competition", func(t *testing.T) {
		done := make(chan bool)

		repo.On("GetById", 564).Return(&model.Competition{}, errors.New("not Found"))
		repo.On("Insert", mock.Anything).Return(nil)
		repo.AssertNotCalled(t, "Update", mock.Anything)
		processor.Process("competition", "", done)
	})

	t.Run("updates existing competition", func(t *testing.T) {
		done := make(chan bool)

		c := newCompetition(1)
		repo.On("GetById", 564).Return(c, nil)
		repo.On("Update", &c).Return(nil)
		repo.MethodCalled("Update", &c)
		repo.AssertNotCalled(t, "Insert", mock.Anything)
		processor.Process("competition", "", done)
	})
}

type mockRepository struct {
	mock.Mock
}

func (m mockRepository) Insert(c *model.Competition) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRepository) Update(c *model.Competition) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockRepository) GetById(id int) (*model.Competition, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Competition)
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

func leagueResponse() sportmonks.LeaguesResponse {
	c := clientLeague()

	m := sportmonks.Meta{}
	m.Pagination.Total = 1
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := sportmonks.LeaguesResponse{}
	res.Data = append(res.Data, c)
	res.Meta = m

	return res
}

func clientLeague() sportmonks.League {
	return sportmonks.League{
		ID:              564,
		LegacyID:        3491,
		CountryID:       32,
		Name:            "Serie A",
		IsCup:           false,
		CurrentSeasonID: 23,
		CurrentRoundID:  98,
		CurrentStageID:  87,
		LiveStandings:   true,
		Coverage: struct {
			TopscorerGoals   bool `json:"topscorer_goals"`
			TopscorerAssists bool `json:"topscorer_assists"`
			TopscorerCards   bool `json:"topscorer_cards"`
		}{
			TopscorerGoals:   true,
			TopscorerAssists: false,
			TopscorerCards:   true,
		},
		Seasons: struct {
			Data []sportmonks.Season `json:"data"`
		}{
			Data: []sportmonks.Season{},
		},
	}
}
