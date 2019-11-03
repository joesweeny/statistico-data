package sportmonks_test

import (
	"bytes"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSquadRequester_SquadsBySeasonIDs(t *testing.T) {
	t.Run("returns a channel containing squad struct", func(t *testing.T) {
		t.Helper()

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(teamSquadResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewSquadRequester(&client, logger)

		ch := requester.SquadsBySeasonIDs([]uint64{uint64(435), uint64(33), uint64(2)})

		x := <- ch
		y := <- ch
		z := <- ch

		a := assert.New(t)
		//
		a.Equal(uint64(1), x.TeamID)
		a.Equal([]uint64{uint64(219591)}, x.PlayerIDs)

		a.Equal(uint64(1), y.TeamID)
		a.Equal([]uint64{uint64(219591)}, x.PlayerIDs)

		a.Equal(uint64(1), z.TeamID)
		a.Equal([]uint64{uint64(219591)}, x.PlayerIDs)
	})
}

var teamSquadResponse = `{
	"data": [
		{
			"id": 1,
			"legacy_id": 377,
			"name": "West Ham United",
			"short_code": "WHU",
			"twitter": "@WestHamUtd",
			"country_id": 462,
			"national_team": false,
			"founded": 1895,
			"logo_path": "https:\/\/cdn.sportmonks.com\/images\/soccer\/teams\/1\/1.png",
			"venue_id": 214,
			"current_season_id": 16036,
			"squad": {
				"data": [
					{
						"player_id": 219591,
						"position_id": 2,
						"number": 4,
						"captain": 0,
						"injured": false,
						"minutes": 90,
						"appearences": 2,
						"lineups": 1,
						"substitute_in": 1,
						"substitute_out": 0,
						"substitutes_on_bench": 7,
						"goals": 0,
						"assists": 0,
						"saves": null,
						"inside_box_saves": null,
						"dispossesed": null,
						"interceptions": null,
						"yellowcards": 1,
						"yellowred": 0,
						"redcards": 0,
						"tackles": null,
						"blocks": null,
						"hit_post": null
					}
				]
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
