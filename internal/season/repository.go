package season

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	Insert(s *model.Season) error
	Update(s *model.Season) error
	Id(id int64) (*model.Season, error)
	Ids() ([]int64, error)
	CurrentSeasonIds() ([]int64, error)
}
