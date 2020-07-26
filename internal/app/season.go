package app

import (
	"time"
)

// Season domain entity.
type Season struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	CompetitionID uint64    `json:"league_id"`
	IsCurrent     bool      `json:"current"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SeasonRepository provides an interface to persist Season domain struct objects to a storage engine.
type SeasonRepository interface {
	Insert(s *Season) error
	Update(s *Season) error
	ByID(id uint64) (*Season, error)
	IDs() ([]uint64, error)
	CurrentSeasonIDs() ([]uint64, error)
	ByCompetitionId(id uint64, sort string) ([]Season, error)
}

// SeasonRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type SeasonRequester interface {
	Seasons() <-chan *Season
}
