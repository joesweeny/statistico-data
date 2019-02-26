package squad

import "github.com/joesweeny/statistico-data/internal/model"

type Repository interface {
	Insert(m *model.Squad) error
	Update(m *model.Squad) error
	BySeasonAndTeam(seasonId, teamId int) (*model.Squad, error)
	All() ([]model.Squad, error)
	CurrentSeason() ([]model.Squad, error)
}
