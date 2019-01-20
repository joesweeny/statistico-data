package season

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(s model.Season) error
	GetById(id int) (model.Season, error)
	GetByLeagueId(id int) (model.Season, error)
}
