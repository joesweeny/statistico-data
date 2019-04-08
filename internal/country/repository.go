package country

import (
	"github.com/statistico/statistico-data/internal/model"
)

type Repository interface {
	Insert(c *model.Country) error
	Update(c *model.Country) error
	GetById(id int) (*model.Country, error)
}
