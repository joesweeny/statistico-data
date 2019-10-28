package app

import (
	"time"
)

type Round struct {
	ID        int64       `json:"id"`
	Name      string    `json:"name"`
	SeasonID  int64       `json:"season_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoundRepository provides an interface to persist Round domain struct objects to a storage engine.
type RoundRepository interface {
	Insert(r *Round) error
	Update(r *Round) error
	ByID(id int64) (*Round, error)
}
