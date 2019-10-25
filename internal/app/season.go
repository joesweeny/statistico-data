package app

import (
	"time"
)

// Season domain entity.
type Season struct {
	ID        int64       `json:"id"`
	Name      string    `json:"name"`
	CompetitionID  int64       `json:"league_id"`
	IsCurrent bool      `json:"current"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SeasonRepository provides an interface to persist Season domain struct objects to a storage engine.
type SeasonRepository interface {
	Insert(s *Season) error
	Update(s *Season) error
	ByID(id int64) (*Season, error)
	IDs() ([]int64, error)
	CurrentSeasonIDs() ([]int64, error)
}
