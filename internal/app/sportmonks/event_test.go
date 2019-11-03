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

func TestEventRequester_EventsByFixtureIDs(t *testing.T) {
	t.Run("returns two channels containing goal and substitution events", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fixtureEventsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewEventRequester(&client, logger)

		goals, subs := requester.EventsByFixtureIDs([]int64{int64(5), int64(23)})

		goalOne := <- goals
		goalTwo := <- goals

		subOne := <- subs
		subTwo := <- subs

		a := assert.New(t)

		a.Equal(int64(11867297001), goalOne.ID)
		a.Equal(int64(78), goalOne.TeamID)
		a.Equal(int64(11867297), goalOne.FixtureID)
		a.Equal(int64(95776), goalOne.PlayerID)
		a.Equal(int64(13452), *goalOne.PlayerAssistID)
		a.Equal(3, goalOne.Minute)
		a.Equal("1-0", goalOne.Score)

		a.Equal(int64(11867297001), goalTwo.ID)
		a.Equal(int64(78), goalTwo.TeamID)
		a.Equal(int64(11867297), goalTwo.FixtureID)
		a.Equal(int64(95776), goalTwo.PlayerID)
		a.Equal(int64(13452), *goalTwo.PlayerAssistID)
		a.Equal(3, goalTwo.Minute)
		a.Equal("1-0", goalTwo.Score)

		a.Equal(int64(11867325003), subOne.ID)
		a.Equal(int64(21), subOne.TeamID)
		a.Equal(int64(11867325), subOne.FixtureID)
		a.Equal(int64(1384), subOne.PlayerInID)
		a.Equal(int64(3530), subOne.PlayerOutID)
		a.Equal(54, subOne.Minute)
		a.Nil(subOne.Injured)

		a.Equal(int64(11867325003), subTwo.ID)
		a.Equal(int64(21), subTwo.TeamID)
		a.Equal(int64(11867325), subTwo.FixtureID)
		a.Equal(int64(1384), subTwo.PlayerInID)
		a.Equal(int64(3530), subTwo.PlayerOutID)
		a.Equal(54, subTwo.Minute)
		a.Nil(subTwo.Injured)
	})
}

var fixtureEventsResponse = `{
	"data": {
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
		},
		"deleted": false,
		"goals": {
			"data": [
				{
					"id": 11867297001,
					"team_id": "78",
					"type": "goal",
					"fixture_id": 11867297,
					"player_id": 95776,
					"player_name": "N. Maupay",
					"player_assist_id": 13452,
					"player_assist_name": null,
					"minute": 3,
					"extra_minute": null,
					"reason": null,
					"result": "1-0"
				}
			]
		},
		"substitutions": {
			"data": [
				{
					  "id": 11867325003,
					  "team_id": "21",
					  "type": "subst",
					  "fixture_id": 11867325,
					  "player_in_id": 1384,
					  "player_in_name": "B. Sharp",
					  "player_out_id": 3530,
					  "player_out_name": "C. Robinson",
					  "minute": 54,
					  "extra_minute": null,
					  "injuried": null
				}
			]
		}
  	}
}`