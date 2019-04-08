package season

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createSeason(s *sportmonks.Season) *model.Season {
	return &model.Season{
		ID:        s.ID,
		Name:      s.Name,
		LeagueID:  s.LeagueID,
		IsCurrent: s.IsCurrentSeason,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f Factory) updateSeason(s *sportmonks.Season, m *model.Season) *model.Season {
	m.Name = s.Name
	m.IsCurrent = s.IsCurrentSeason
	m.UpdatedAt = f.Clock.Now()

	return m
}
