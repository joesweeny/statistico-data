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

func TestResultRequester_ResultsBySeasonIDs(t *testing.T) {
	t.Run("returns result struct channel", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(seasonResultsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewResultRequester(&client, logger)

		ch := requester.ResultsBySeasonIDs([]uint64{11867285})

		result := <-ch

		a := assert.New(t)

		a.Equal(uint64(11867285), result.FixtureID)
		a.Equal("Good", *result.PitchCondition)
		a.Equal("4-1-4-1", *result.HomeFormation)
		a.Equal("4-2-3-1", *result.AwayFormation)
		a.Equal(2, *result.HomeScore)
		a.Equal(0, *result.AwayScore)
		a.Nil(result.HomePenScore)
		a.Nil(result.AwayPenScore)
		a.Equal("1-0", *result.HalfTimeScore)
		a.Equal("2-0", *result.FullTimeScore)
		a.Nil(result.ExtraTimeScore)
		a.Equal(11, *result.HomeLeaguePosition)
		a.Equal(6, *result.AwayLeaguePosition)
		a.Equal(90, *result.Minutes)
		a.Nil(result.AddedTime)
		a.Nil(result.ExtraTime)
		a.Nil(result.InjuryTime)
	})
}

var seasonResultsResponse = `{
	"data": {
		"id": 16029,
		"name": "2019/2020",
		"league_id": 2,
		"is_current_season": true,
		"current_round_id": 183973,
		"current_stage_id": 77443828,
		"results": {
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
					"pitch": "Good",
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
	}
}`
