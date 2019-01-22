package round

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(round *model.Round) error
	Update(round *model.Round) error
	GetById(id int) (*model.Round, error)
}
