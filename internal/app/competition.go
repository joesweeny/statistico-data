package app

import (
	"time"
)

// Competition domain entity.
type Competition struct {
	ID        uint64       `json:"id"`
	Name      string    `json:"name"`
	CountryID uint64       `json:"country_id"`
	IsCup     bool      `json:"is_cup"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CompetitionRepository provides an interface to persist Competition domain struct objects to a storage engine.
type CompetitionRepository interface {
	Insert(c *Competition) error
	Update(c *Competition) error
	ByID(id uint64) (*Competition, error)
}

// CompetitionRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type CompetitionRequester interface {
	Competitions() <-chan *Competition
}
