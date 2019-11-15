package app

import (
	"time"
)

// Fixture domain entity.
type Fixture struct {
	ID         uint64      `json:"id"`
	SeasonID   uint64       `json:"season_id"`
	RoundID    *uint64      `json:"round_id"`
	VenueID    *uint64      `json:"venue_id"`
	HomeTeamID uint64       `json:"home_team_id"`
	AwayTeamID uint64       `json:"away_team_id"`
	RefereeID  *uint64      `json:"referee_id"`
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FixtureRepository provides an interface to persist Fixture domain struct objects to a storage engine.
type FixtureRepository interface {
	Insert(f *Fixture) error
	Update(f *Fixture) error
	ByID(id uint64) (*Fixture, error)
	ByTeamID(id uint64, limit int32, before time.Time) ([]Fixture, error)
	Get(q FixtureRepositoryQuery) ([]Fixture, error)
	GetIDs(q FixtureRepositoryQuery) ([]uint64, error)
}

type FixtureRepositoryQuery struct {
	SeasonID *uint64
	HomeTeamID *uint64
	AwayTeamID *uint64
	DateFrom *time.Time
	DateTo *time.Time
	Limit *uint64
}

// FixtureRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type FixtureRequester interface {
	FixturesBySeasonIDs(ids []uint64) <-chan *Fixture
}
