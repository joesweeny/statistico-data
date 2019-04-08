package result

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	Insert(r *model.Result) error
	Update(r *model.Result) error
	GetByFixtureId(id int) (*model.Result, error)
}
