package model

import "time"

type Result struct {
	FixtureID          int       `json:"fixture_id"`
	PitchCondition     *string   `json:"pitch_condition"`
	HomeFormation      *string   `json:"home_formation"`
	AwayFormation      *string   `json:"away_formation"`
	HomeScore          *int      `json:"home_score"`
	AwayScore          *int      `json:"away_score"`
	HomePenScore       *int      `json:"home_pen_score"`
	AwayPenScore       *int      `json:"away_pen_score"`
	HalfTimeScore      *string   `json:"half_time_score"`
	FullTimeScore      *string   `json:"full_time_score"`
	ExtraTimeScore     *string   `json:"extra_time_score"`
	HomeLeaguePosition *int      `json:"home_league_position"`
	AwayLeaguePosition *int      `json:"away_league_position"`
	Minutes            *int      `json:"minutes"`
	Seconds            *int      `json:"seconds"`
	AddedTime          *int      `json:"added_time"`
	ExtraTime          *int      `json:"extra_time"`
	InjuryTime         *int      `json:"injury_time"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
