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
	"time"
)

func TestTeamStatsRequester_TeamStatsByFixtureIDs(t *testing.T) {
	t.Run("returns team stats struct channel", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fixtureTeamStatsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewTeamStatsRequester(&client, logger)

		ch := requester.TeamStatsByFixtureIDs([]uint64{11867285})

		home := <-ch
		away := <-ch

		a := assert.New(t)

		a.Equal(uint64(830), home.TeamID)
		a.Equal(uint64(11966141), home.FixtureID)
		a.Equal(14, *home.TeamShots.Total)
		a.Equal(4, *home.TeamShots.OnGoal)
		a.Equal(10, *home.TeamShots.OffGoal)
		a.Equal(6, *home.TeamShots.InsideBox)
		a.Equal(9, *home.TeamShots.OutsideBox)
		a.Nil(home.TeamShots.Blocked)
		a.Equal(313, *home.TeamPasses.Total)
		a.Equal(230, *home.TeamPasses.Accuracy)
		a.Equal(float32(42.5), *home.TeamPasses.Percentage)
		a.Equal(104, *home.TeamAttacks.Total)
		a.Equal(52, *home.TeamAttacks.Dangerous)
		a.Equal(13, *home.Fouls)
		a.Equal(8, *home.Corners)
		a.Equal(1, *home.Offsides)
		a.Equal(35, *home.Possession)
		a.Equal(1, *home.YellowCards)
		a.Equal(0, *home.RedCards)
		a.Equal(0, *home.Saves)
		a.Equal(3, *home.Substitutions)
		a.Equal(11, *home.GoalKicks)
		a.Equal(14, *home.GoalAttempts)
		a.Equal(16, *home.FreeKicks)
		a.Equal(30, *home.ThrowIns)

		a.Equal(uint64(19), away.TeamID)
		a.Equal(uint64(11966141), away.FixtureID)
		a.Equal(7, *away.TeamShots.Total)
		a.Equal(1, *away.TeamShots.OnGoal)
		a.Equal(6, *away.TeamShots.OffGoal)
		a.Equal(4, *away.TeamShots.InsideBox)
		a.Equal(3, *away.TeamShots.OutsideBox)
		a.Nil(away.TeamShots.Blocked)
		a.Equal(588, *away.TeamPasses.Total)
		a.Equal(519, *away.TeamPasses.Accuracy)
		a.Equal(float32(0), *away.TeamPasses.Percentage)
		a.Equal(148, *away.TeamAttacks.Total)
		a.Equal(73, *away.TeamAttacks.Dangerous)
		a.Equal(12, *away.Fouls)
		a.Equal(6, *away.Corners)
		a.Equal(4, *away.Offsides)
		a.Equal(65, *away.Possession)
		a.Equal(2, *away.YellowCards)
		a.Equal(0, *away.RedCards)
		a.Equal(3, *away.Saves)
		a.Equal(3, *away.Substitutions)
		a.Equal(12, *away.GoalKicks)
		a.Equal(6, *away.GoalAttempts)
		a.Equal(14, *away.FreeKicks)
		a.Equal(27, *away.ThrowIns)
	})
}

func TestTeamStatsRequester_TeamStatsByDate(t *testing.T) {
	t.Run("returns team stats struct channel", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fixturesTeamStatsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewTeamStatsRequester(&client, logger)

		ch := requester.TeamStatsByDate(time.Now(), []uint64{11867285})

		home := <-ch
		away := <-ch

		a := assert.New(t)

		a.Equal(uint64(830), home.TeamID)
		a.Equal(uint64(11966141), home.FixtureID)
		a.Equal(14, *home.TeamShots.Total)
		a.Equal(4, *home.TeamShots.OnGoal)
		a.Equal(10, *home.TeamShots.OffGoal)
		a.Equal(6, *home.TeamShots.InsideBox)
		a.Equal(9, *home.TeamShots.OutsideBox)
		a.Nil(home.TeamShots.Blocked)
		a.Equal(313, *home.TeamPasses.Total)
		a.Equal(230, *home.TeamPasses.Accuracy)
		a.Equal(float32(42.5), *home.TeamPasses.Percentage)
		a.Equal(104, *home.TeamAttacks.Total)
		a.Equal(52, *home.TeamAttacks.Dangerous)
		a.Equal(13, *home.Fouls)
		a.Equal(8, *home.Corners)
		a.Equal(1, *home.Offsides)
		a.Equal(35, *home.Possession)
		a.Equal(1, *home.YellowCards)
		a.Equal(0, *home.RedCards)
		a.Equal(0, *home.Saves)
		a.Equal(3, *home.Substitutions)
		a.Equal(11, *home.GoalKicks)
		a.Equal(14, *home.GoalAttempts)
		a.Equal(16, *home.FreeKicks)
		a.Equal(30, *home.ThrowIns)

		a.Equal(uint64(19), away.TeamID)
		a.Equal(uint64(11966141), away.FixtureID)
		a.Equal(7, *away.TeamShots.Total)
		a.Equal(1, *away.TeamShots.OnGoal)
		a.Equal(6, *away.TeamShots.OffGoal)
		a.Equal(4, *away.TeamShots.InsideBox)
		a.Equal(3, *away.TeamShots.OutsideBox)
		a.Nil(away.TeamShots.Blocked)
		a.Equal(588, *away.TeamPasses.Total)
		a.Equal(519, *away.TeamPasses.Accuracy)
		a.Equal(float32(0), *away.TeamPasses.Percentage)
		a.Equal(148, *away.TeamAttacks.Total)
		a.Equal(73, *away.TeamAttacks.Dangerous)
		a.Equal(12, *away.Fouls)
		a.Equal(6, *away.Corners)
		a.Equal(4, *away.Offsides)
		a.Equal(65, *away.Possession)
		a.Equal(2, *away.YellowCards)
		a.Equal(0, *away.RedCards)
		a.Equal(3, *away.Saves)
		a.Equal(3, *away.Substitutions)
		a.Equal(12, *away.GoalKicks)
		a.Equal(6, *away.GoalAttempts)
		a.Equal(14, *away.FreeKicks)
		a.Equal(27, *away.ThrowIns)
	})
}

var fixtureTeamStatsResponse = `{
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
		},
		"stats": {
		  	"data": [
				{
					"team_id": 830,
					"fixture_id": 11966141,
					"shots": {
						"total": 14,
						"ongoal": 4,
						"offgoal": 10,
						"insidebox": 6,
						"outsidebox": 9
					},
					"passes": {
						"total": 313,
						"accurate": 230,
						"percentage": 42.5
					},
					"attacks": {
						"attacks": 104,
						"dangerous_attacks": 52
					},
					"fouls": 13,
					"corners": 8,
					"offsides": 1,
					"possessiontime": 35,
					"yellowcards": 1,
					"redcards": 0,
					"yellowredcards": 0,
					"saves": 0,
					"substitutions": 3,
					"goal_kick": 11,
					"goal_attempts": 14,
					"free_kick": 16,
					"throw_in": 30,
					"ball_safe": 77,
					"goals": 1,
					"penalties": 0,
					"injuries": 1
					},
					{
					"team_id": 19,
					"fixture_id": 11966141,
					"shots": {
						"total": 7,
						"ongoal": 1,
						"offgoal": 6,
						"insidebox": 4,
						"outsidebox": 3
					},
					"passes": {
						"total": 588,
						"accurate": 519,
						"percentage": 0
					},
					"attacks": {
						"attacks": 148,
						"dangerous_attacks": 73
					},
					"fouls": 12,
					"corners": 6,
					"offsides": 4,
					"possessiontime": 65,
					"yellowcards": 2,
					"redcards": 0,
					"yellowredcards": 0,
					"saves": 3,
					"substitutions": 3,
					"goal_kick": 12,
					"goal_attempts": 6,
					"free_kick": 14,
					"throw_in": 27,
					"ball_safe": 99,
					"goals": 1,
					"penalties": 0,
					"injuries": 0
				}
		  	]
		}
	}
}`

var fixturesTeamStatsResponse = `{
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
			},
			"stats": {
				"data": [
					{
						"team_id": 830,
						"fixture_id": 11966141,
						"shots": {
							"total": 14,
							"ongoal": 4,
							"offgoal": 10,
							"insidebox": 6,
							"outsidebox": 9
						},
						"passes": {
							"total": 313,
							"accurate": 230,
							"percentage": 42.5
						},
						"attacks": {
							"attacks": 104,
							"dangerous_attacks": 52
						},
						"fouls": 13,
						"corners": 8,
						"offsides": 1,
						"possessiontime": 35,
						"yellowcards": 1,
						"redcards": 0,
						"yellowredcards": 0,
						"saves": 0,
						"substitutions": 3,
						"goal_kick": 11,
						"goal_attempts": 14,
						"free_kick": 16,
						"throw_in": 30,
						"ball_safe": 77,
						"goals": 1,
						"penalties": 0,
						"injuries": 1
						},
						{
						"team_id": 19,
						"fixture_id": 11966141,
						"shots": {
							"total": 7,
							"ongoal": 1,
							"offgoal": 6,
							"insidebox": 4,
							"outsidebox": 3
						},
						"passes": {
							"total": 588,
							"accurate": 519,
							"percentage": 0
						},
						"attacks": {
							"attacks": 148,
							"dangerous_attacks": 73
						},
						"fouls": 12,
						"corners": 6,
						"offsides": 4,
						"possessiontime": 65,
						"yellowcards": 2,
						"redcards": 0,
						"yellowredcards": 0,
						"saves": 3,
						"substitutions": 3,
						"goal_kick": 12,
						"goal_attempts": 6,
						"free_kick": 14,
						"throw_in": 27,
						"ball_safe": 99,
						"goals": 1,
						"penalties": 0,
						"injuries": 0
					}
				]
			}
		}
	]
}`
