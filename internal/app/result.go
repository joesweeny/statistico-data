package app

import (
	"time"
)

// Result domain entity.
type Result struct {
	FixtureID          uint64    `json:"fixture_id"`
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
	AddedTime          *int      `json:"added_time"`
	ExtraTime          *int      `json:"extra_time"`
	InjuryTime         *int      `json:"injury_time"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ResultRepository provides an interface to persist Result domain struct objects to a storage engine.
type ResultRepository interface {
	Insert(r *Result) error
	Update(r *Result) error
	ByFixtureID(id uint64) (*Result, error)
}

// ResultRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type ResultRequester interface {
	ResultsBySeasonIDs(seasonIDs []uint64) <-chan Result
}
