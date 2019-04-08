package competition

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createCompetition(s *sportmonks.League) *model.Competition {
	return &model.Competition{
		ID:        s.ID,
		Name:      s.Name,
		CountryID: s.CountryID,
		IsCup:     s.IsCup,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f Factory) updateCompetition(s *sportmonks.League, m *model.Competition) *model.Competition {
	m.Name = s.Name
	m.IsCup = s.IsCup
	m.UpdatedAt = f.Clock.Now()

	return m
}
