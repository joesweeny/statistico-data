package round

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(r *model.Round) error
	Update(r *model.Round) error
	GetById(id int) (*model.Round, error)
}
