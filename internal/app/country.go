package app

import (
	"time"
)

// Country domain entity.
type Country struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Continent string    `json:"continent"`
	ISO       string    `json:"iso"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CountryRepository provides an interface to persist Country domain struct objects to a storage engine.
type CountryRepository interface {
	Insert(c *Country) error
	Update(c *Country) error
	GetById(id int64) (*Country, error)
}

// CountryRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type CountryRequester interface {
	Countries() <-chan *Country
}
