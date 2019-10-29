package app

import (
	"time"
)

// Team domain entity.
type Team struct {
	ID           int64      `json:"id"`
	Name         string    `json:"name"`
	ShortCode    *string   `json:"short_code"`
	CountryID    *int64      `json:"country_id"`
	VenueID      int64       `json:"venue_id"`
	NationalTeam bool      `json:"national_team"`
	Founded      *int      `json:"founded"`
	Logo         *string   `json:"logo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TeamRepository provides an interface to persist Team domain struct objects to a storage engine.
type TeamRepository interface {
	Insert(t *Team) error
	Update(t *Team) error
	ByID(id int64) (*Team, error)
}

// TeamRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type TeamRequester interface {
	TeamsBySeasonID(seasonID int64) <-chan *Team
}
