package app

import (
	"time"
)

// Player domain entity.
type Player struct {
	ID          int64       `json:"id"`
	CountryId   int64       `json:"country_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	BirthPlace  *string   `json:"birth_place"`
	DateOfBirth *string   `json:"date_of_birth"`
	PositionID  int       `json:"position_id"`
	Image       string   `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlayerRepository provides an interface to persist Player domain struct objects to a storage engine.
type PlayerRepository interface {
	Insert(m *Player) error
	Update(m *Player) error
	ByID(id int64) (*Player, error)
}

// PlayerRequester provides an interface allowing this application to request player data from an external
// data provider
type PlayerRequester interface {
	PlayerByID(id int64) (*Player, error)
}
