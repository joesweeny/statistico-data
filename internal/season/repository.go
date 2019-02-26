package season

import "github.com/joesweeny/statistico-data/internal/model"

type Repository interface {
	Insert(s *model.Season) error
	Update(s *model.Season) error
	Id(id int) (*model.Season, error)
	Ids() ([]int, error)
	CurrentSeasonIds() ([]int, error)
}
