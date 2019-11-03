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

func TestFixtureRequester_FixturesBySeasonIds(t *testing.T) {
	t.Run("returns a channel of fixture struct", func(t *testing.T) {
		t.Helper()

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(seasonFixturesResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewFixtureRequester(&client, logger)

		ch := requester.FixturesBySeasonIDs([]uint64{uint64(435), uint64(33), uint64(2)})

		x := <-ch
		y := <-ch
		z := <-ch

		a := assert.New(t)

		a.Equal(uint64(11867285), x.ID)
		a.Equal(uint64(16036), x.SeasonID)
		a.Equal(uint64(169657), *x.RoundID)
		a.Equal(uint64(214), *x.VenueID)
		a.Equal(uint64(1), x.HomeTeamID)
		a.Equal(uint64(14), x.AwayTeamID)
		a.Equal(uint64(14532), *x.RefereeID)
		a.Equal("2019-09-22 13:00:00 +0000 UTC", x.Date.String())

		a.Equal(uint64(11867285), y.ID)
		a.Equal(uint64(16036), y.SeasonID)
		a.Equal(uint64(169657), *y.RoundID)
		a.Equal(uint64(214), *y.VenueID)
		a.Equal(uint64(1), y.HomeTeamID)
		a.Equal(uint64(14), y.AwayTeamID)
		a.Equal(uint64(14532), *y.RefereeID)
		a.Equal("2019-09-22 13:00:00 +0000 UTC", y.Date.String())

		a.Equal(uint64(11867285), z.ID)
		a.Equal(uint64(16036), z.SeasonID)
		a.Equal(uint64(169657), *z.RoundID)
		a.Equal(uint64(214), *z.VenueID)
		a.Equal(uint64(1), z.HomeTeamID)
		a.Equal(uint64(14), z.AwayTeamID)
		a.Equal(uint64(14532), *z.RefereeID)
		a.Equal("2019-09-22 13:00:00 +0000 UTC", z.Date.String())
	})
}

var seasonFixturesResponse = `{
	"data": {
		"id": 16029,
		"name": "2019/2020",
		"league_id": 2,
		"is_current_season": true,
		"current_round_id": 183973,
		"current_stage_id": 77443828,
		"fixtures": {
			"data": [
				{
					"id": 11867285,
					"league_id": 8,
					"season_id": 16036,
					"stage_id": 77443862,
					"round_id": 169657,
					"group_id": null,
					"aggregate_id": null,
					"venue_id": 214,
					"referee_id": 14532,
					"localteam_id": 1,
					"visitorteam_id": 14,
					"winner_team_id": 1,
					"weather_report": {
					  "code": "rain",
					  "type": "shower rain",
					  "icon": "https:\/\/cdn.sportmonks.com\/images\/weather\/09d.png",
					  "temperature": {
						"temp": 62.96,
						"unit": "fahrenheit"
					  },
					  "temperature_celcius": {
						"temp": 17.2,
						"unit": "celcius"
					  },
					  "clouds": "75%",
					  "humidity": "82%",
					  "pressure": 1004,
					  "wind": {
						"speed": "5.82 m\/s",
						"degree": 240
					  },
					  "coordinates": {
						"lat": 51.51,
						"lon": -0.13
					  },
					  "updated_at": "2019-09-22T14:45:05.289505Z"
					},
					"commentaries": true,
					"attendance": 59936,
					"pitch": null,
					"details": null,
					"neutral_venue": false,
					"winning_odds_calculated": true,
					"formations": {
					  "localteam_formation": "4-1-4-1",
					  "visitorteam_formation": "4-2-3-1"
					},
					"scores": {
					  "localteam_score": 2,
					  "visitorteam_score": 0,
					  "localteam_pen_score": null,
					  "visitorteam_pen_score": null,
					  "ht_score": "1-0",
					  "ft_score": "2-0",
					  "et_score": null,
					  "ps_score": null
					},
					"time": {
					  "status": "FT",
					  "starting_at": {
						"date_time": "2019-09-22 13:00:00",
						"date": "2019-09-22",
						"time": "13:00:00",
						"timestamp": 1569157200,
						"timezone": "UTC"
					  },
					  "minute": 90,
					  "second": null,
					  "added_time": null,
					  "extra_minute": null,
					  "injury_time": null
					},
					"coaches": {
					  "localteam_coach_id": 523898,
					  "visitorteam_coach_id": 524307
					},
					"standings": {
					  "localteam_position": 11,
					  "visitorteam_position": 6
					},
					"assistants": {
					  "first_assistant_id": 12794,
					  "second_assistant_id": 12798,
					  "fourth_official_id": 15270
					},
					"leg": "1\/1",
					"colors": {
					  "localteam": {
						"color": "#832034",
						"kit_colors": "#C0D6FE,#C0D6FE,#832034,#832034,#999999,#888888,#832034"
					  },
					  "visitorteam": {
						"color": null,
						"kit_colors": null
					  }
					}
				}
			]
		}
	},
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
