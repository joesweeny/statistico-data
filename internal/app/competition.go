package app

import (
	"time"
)

// Competition domain entity.
type Competition struct {
	ID        int64       `json:"id"`
	Name      string    `json:"name"`
	CountryID int64       `json:"country_id"`
	IsCup     bool      `json:"is_cup"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CompetitionRepository provides an interface to persist Competition domain struct objects to a storage engine.
type CompetitionRepository interface {
	Insert(c *Competition) error
	Update(c *Competition) error
	ByID(id int64) (*Competition, error)
}
