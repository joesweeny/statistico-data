package statistico

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

// CountryDataService provides an interface allowing this application to retrieve data from an external
// data provider and filtering through the channel provided as the only argument
type CountryDataService interface {
	Countries(ch chan<- *Country)
}
