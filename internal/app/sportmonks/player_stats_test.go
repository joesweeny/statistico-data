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

func TestPlayerStatsRequester_PlayerStatsByFixtureIDs(t *testing.T) {
	t.Run("returns player stats struct channel", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fixturePlayerStatsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewPlayerStatsRequester(&client, logger)

		ch := requester.PlayerStatsByFixtureIDs([]uint64{11867285})

		pl := <-ch
		a := assert.New(t)

		a.Equal(uint64(11937554), pl.FixtureID)
		a.Equal(uint64(129665), pl.PlayerID)
		a.Equal(uint64(613), pl.TeamID)
		a.Equal("M", *pl.Position)
		a.Equal(8, *pl.FormationPosition)
		a.Equal(false, pl.IsSubstitute)
		a.Equal(0, *pl.PlayerShots.Total)
		a.Equal(0, *pl.PlayerShots.OnGoal)
		a.Equal(0, *pl.PlayerGoals.Scored)
		a.Equal(0, *pl.PlayerGoals.Conceded)
		a.Equal(0, *pl.PlayerFouls.Drawn)
		a.Equal(2, *pl.PlayerFouls.Committed)
		a.Equal(1, *pl.YellowCards)
		a.Equal(0, *pl.RedCard)
		a.Nil(pl.PlayerPenalties.Scored)
		a.Nil(pl.PlayerPenalties.Missed)
		a.Nil(pl.PlayerPenalties.Saved)
		a.Equal(0, *pl.PlayerPenalties.Committed)
		a.Equal(0, *pl.PlayerPenalties.Won)
		a.Equal(2, *pl.PlayerCrosses.Total)
		a.Equal(1, *pl.PlayerCrosses.Accuracy)
		a.Equal(25, *pl.PlayerPasses.Total)
		a.Equal(89, *pl.PlayerPasses.Accuracy)
		a.Equal(0, *pl.Assists)
		a.Equal(0, *pl.Offsides)
		a.Equal(0, *pl.Saves)
		a.Equal(0, *pl.HitWoodwork)
		a.Equal(1, *pl.Tackles)
		a.Equal(0, *pl.Blocks)
		a.Equal(0, *pl.Interceptions)
		a.Equal(0, *pl.Clearances)
		a.Equal(85, *pl.MinutesPlayed)
	})
}

func TestPlayerStatsRequester_PlayerStatsByDate(t *testing.T) {
	t.Run("returns players stats struct channel", func(t *testing.T) {
		t.Helper()

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fixturesPlayerStatsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewPlayerStatsRequester(&client, logger)

		ch := requester.PlayerStatsByDate(time.Now(), []uint64{11867285})

		pl := <-ch
		a := assert.New(t)

		a.Equal(uint64(11937554), pl.FixtureID)
		a.Equal(uint64(129665), pl.PlayerID)
		a.Equal(uint64(613), pl.TeamID)
		a.Equal("M", *pl.Position)
		a.Equal(8, *pl.FormationPosition)
		a.Equal(false, pl.IsSubstitute)
		a.Equal(0, *pl.PlayerShots.Total)
		a.Equal(0, *pl.PlayerShots.OnGoal)
		a.Equal(0, *pl.PlayerGoals.Scored)
		a.Equal(0, *pl.PlayerGoals.Conceded)
		a.Equal(0, *pl.PlayerFouls.Drawn)
		a.Equal(2, *pl.PlayerFouls.Committed)
		a.Equal(1, *pl.YellowCards)
		a.Equal(0, *pl.RedCard)
		a.Nil(pl.PlayerPenalties.Scored)
		a.Nil(pl.PlayerPenalties.Missed)
		a.Nil(pl.PlayerPenalties.Saved)
		a.Equal(0, *pl.PlayerPenalties.Committed)
		a.Equal(0, *pl.PlayerPenalties.Won)
		a.Equal(2, *pl.PlayerCrosses.Total)
		a.Equal(1, *pl.PlayerCrosses.Accuracy)
		a.Equal(25, *pl.PlayerPasses.Total)
		a.Equal(89, *pl.PlayerPasses.Accuracy)
		a.Equal(0, *pl.Assists)
		a.Equal(0, *pl.Offsides)
		a.Equal(0, *pl.Saves)
		a.Equal(0, *pl.HitWoodwork)
		a.Equal(1, *pl.Tackles)
		a.Equal(0, *pl.Blocks)
		a.Equal(0, *pl.Interceptions)
		a.Equal(0, *pl.Clearances)
		a.Equal(85, *pl.MinutesPlayed)
	})
}

var fixturePlayerStatsResponse = `{
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
		"lineup": {
		  	"data": [
				 {
					  "team_id": 613,
					  "fixture_id": 11937554,
					  "player_id": 129665,
					  "player_name": "Daniele Baselli",
					  "number": 8,
					  "position": "M",
					  "additional_position": null,
					  "formation_position": 8,
					  "posx": 3,
					  "posy": 2,
					  "captain": false,
					  "stats": {
						"shots": {
						  "shots_total": 0,
						  "shots_on_goal": 0
						},
						"goals": {
						  "scored": 0,
						  "assists": 0,
						  "conceded": 0
						},
						"fouls": {
						  "drawn": 0,
						  "committed": 2
						},
						"cards": {
						  "yellowcards": 1,
						  "redcards": 0,
						  "yellowredcards": 0
						},
						"passing": {
						  "total_crosses": 2,
						  "crosses_accuracy": 1,
						  "passes": 25,
						  "passes_accuracy": 89,
						  "key_passes": 2
						},
						"dribbles": {
						  "attempts": 0,
						  "success": 0,
						  "dribbled_past": 1
						},
						"duels": {
						  "total": 0,
						  "won": 0
						},
						"other": {
						  "offsides": 0,
						  "saves": 0,
						  "inside_box_saves": 0,
						  "pen_scored": null,
						  "pen_missed": null,
						  "pen_saved": null,
						  "pen_committed": 0,
						  "pen_won": 0,
						  "hit_woodwork": 0,
						  "tackles": 1,
						  "blocks": 0,
						  "interceptions": 0,
						  "clearances": 0,
						  "dispossesed": 0,
						  "minutes_played": 85
						}
					}
				}
		  	]
		}
	}
}`

var fixturesPlayerStatsResponse = `{
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
			"lineup": {
				"data": [
					 {
						  "team_id": 613,
						  "fixture_id": 11937554,
						  "player_id": 129665,
						  "player_name": "Daniele Baselli",
						  "number": 8,
						  "position": "M",
						  "additional_position": null,
						  "formation_position": 8,
						  "posx": 3,
						  "posy": 2,
						  "captain": false,
						  "stats": {
							"shots": {
							  "shots_total": 0,
							  "shots_on_goal": 0
							},
							"goals": {
							  "scored": 0,
							  "assists": 0,
							  "conceded": 0
							},
							"fouls": {
							  "drawn": 0,
							  "committed": 2
							},
							"cards": {
							  "yellowcards": 1,
							  "redcards": 0,
							  "yellowredcards": 0
							},
							"passing": {
							  "total_crosses": 2,
							  "crosses_accuracy": 1,
							  "passes": 25,
							  "passes_accuracy": 89,
							  "key_passes": 2
							},
							"dribbles": {
							  "attempts": 0,
							  "success": 0,
							  "dribbled_past": 1
							},
							"duels": {
							  "total": 0,
							  "won": 0
							},
							"other": {
							  "offsides": 0,
							  "saves": 0,
							  "inside_box_saves": 0,
							  "pen_scored": null,
							  "pen_missed": null,
							  "pen_saved": null,
							  "pen_committed": 0,
							  "pen_won": 0,
							  "hit_woodwork": 0,
							  "tackles": 1,
							  "blocks": 0,
							  "interceptions": 0,
							  "clearances": 0,
							  "dispossesed": 0,
							  "minutes_played": 85
							}
						}
					}
				]
			}
		}
	]
}`