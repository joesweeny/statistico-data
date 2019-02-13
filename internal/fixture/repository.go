package fixture

import (
	"github.com/joesweeny/statshub/internal/model"
	"time"
)

type Repository interface {
	Insert(f *model.Fixture) error
	Update(f *model.Fixture) error
	GetById(id int) (*model.Fixture, error)
	Ids() ([]int, error)
	IdsBetween(from, to time.Time) ([]int, error)
}
