package fixture

import (
	"github.com/joesweeny/statistico-data/internal/model"
	"time"
)

type Repository interface {
	Insert(f *model.Fixture) error
	Update(f *model.Fixture) error
	ById(id int) (*model.Fixture, error)
	Ids() ([]int, error)
	IdsBetween(from, to time.Time) ([]int, error)
	Between(from, to time.Time) ([]model.Fixture, error)
	// Id of the Team concerned
	// Limit parameter to limit the number of Fixture structs returned
	// Date constraint returning fixtures from before that date
	ByTeamId(id int64, limit int32, before time.Time) ([]model.Fixture, error)
	// Id of the Season
	// Date constraint returning fixtures from before the given date
	// Order fixture date 'ASC' or 'DESC'
	BySeasonId(id int64, limit int32, before time.Time, order string) ([]model.Fixture, error)
}
