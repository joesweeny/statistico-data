package venue

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(v *model.Venue) error
	Update(v *model.Venue) error
	GetById(id int) (*model.Venue, error)
}
