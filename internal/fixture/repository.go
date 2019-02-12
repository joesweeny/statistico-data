package fixture

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(f *model.Fixture) error
	Update(f *model.Fixture) error
	GetById(id int) (*model.Fixture, error)
	Ids() ([]int, error)
}
