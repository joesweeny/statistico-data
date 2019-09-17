package app

import (
	"time"
)

// Country domain entity
type Country struct {
	ID        int       `json:"id"`
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
	GetById(id int) (*Country, error)
}

// CountryRequester provides an interface allowing this application to request data from an external
// data provider and filtering through the channel provided as the only argument. The requester implementation
// is responsible for closing the channel once successful execution is complete
type CountryRequester interface {
	Countries(ch chan<- *Country)
}
