package result

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(r *model.Result) error
	Update(r *model.Result) error
	GetByFixtureId(id int) (*model.Result, error)
}