package sportmonks_test

import (
	"bytes"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-football-data/internal/app/mock"
	"github.com/statistico/statistico-football-data/internal/app/sportmonks"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCompetitionRequester_Competitions(t *testing.T) {
	t.Run("returns a channel containing competition struct", func(t *testing.T) {
		t.Helper()

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(competitionsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewCompetitionRequester(&client, logger)

		ch := requester.Competitions()

		bun := <-ch

		a := assert.New(t)

		a.Equal(uint64(82), bun.ID)
		a.Equal("Bundesliga", bun.Name)
		a.Equal(uint64(11), bun.CountryID)
		a.Equal(false, bun.IsCup)
	})
}

var competitionsResponse = `{
	"data": [
		{
			"id": 82,
			"active": true,
			"type": "domestic",
			"legacy_id": 4,
			"country_id": 11,
			"logo_path": "https:\/\/cdn.sportmonks.com\/images\/soccer\/leagues\/82.png",
			"name": "Bundesliga",
			"is_cup": false,
			"current_season_id": 16264,
			"current_round_id": 174546,
			"current_stage_id": 77444845,
			"live_standings": true,
			"coverage": {
			  "predictions": true,
			  "topscorer_goals": true,
			  "topscorer_assists": true,
			  "topscorer_cards": true
			}
		}
   ],
	"meta": {
		"pagination": {
			"total": 1,
			"count": 1,
			"per_page": 100,
			"current_page": 1,
			"total_pages": 1,
			"links": {}
		}
	}
}`
