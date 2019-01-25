package season

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(s *model.Season) error
	Update(s *model.Season) error
	GetById(id int) (*model.Season, error)
	GetIds() ([]int, error)
	GetCurrentSeasons() ([]model.Season, error)
}
