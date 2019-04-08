package competition

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	Insert(c *model.Competition) error
	Update(c *model.Competition) error
	GetById(id int) (*model.Competition, error)
}
