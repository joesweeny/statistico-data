package app

import (
	"time"
)

// Round domain entity.
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

// RoundRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type RoundRequester interface {
	RoundsBySeasonIDs(seasonIDs []int64) <-chan *Round
}
