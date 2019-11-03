package app

import "time"

// Venue domain entity.
type Venue struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Surface   *string   `json:"surface"`
	Address   *string   `json:"address"`
	City      *string   `json:"city"`
	Capacity  *int      `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VenueRepository provides an interface to persist Venue domain struct objects to a storage engine.
type VenueRepository interface {
	Insert(v *Venue) error
	Update(v *Venue) error
	GetById(id uint64) (*Venue, error)
}

// VenueRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type VenueRequester interface {
	VenuesBySeasonIDs(seasonIDs []uint64) <-chan *Venue
}
