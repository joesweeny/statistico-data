package country

import (
	"github.com/joesweeny/statshub/internal/model"
)

type Repository interface {
	Insert(c *model.Country) error
	Update(c *model.Country) error
	GetById(id int) (*model.Country, error)
}