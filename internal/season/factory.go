package season

import (
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createSeason(s *sportmonks.Season) *model.Season {
	return &model.Season{
		ID:        s.ID,
		Name:      s.Name,
		LeagueID:  s.LeagueID,
		IsCurrent: s.CurrentSeason,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f Factory) updateSeason(s *sportmonks.Season, m *model.Season) *model.Season {
	m.Name = s.Name
	m.IsCurrent = s.CurrentSeason
	m.UpdatedAt = f.Clock.Now()

	return m
}
