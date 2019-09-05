package statistico

import (
	"time"
)

type Country struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Continent string    `json:"continent"`
	ISO       string    `json:"iso"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Insert(c *Country) error
	Update(c *Country) error
	GetById(id int) (*Country, error)
}
