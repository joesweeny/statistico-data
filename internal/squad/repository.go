package squad

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(m *model.Squad) error
	Update(m *model.Squad) error
	BySeasonAndTeam(seasonId, teamId int) (*model.Squad, error)
	All() (*[]model.Squad, error)
}
