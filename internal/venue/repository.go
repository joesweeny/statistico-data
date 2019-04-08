package venue

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	Insert(v *model.Venue) error
	Update(v *model.Venue) error
	GetById(id int) (*model.Venue, error)
}
