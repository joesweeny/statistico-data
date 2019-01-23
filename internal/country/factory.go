package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createCountry(s *sportmonks.Country) *model.Country {
	return &model.Country{
		ID:        s.ID,
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f Factory) updateCountry(s *sportmonks.Country, m *model.Country) *model.Country {
	m.Name = s.Name
	m.Continent = s.Extra.Continent
	m.ISO = s.Extra.ISO
	m.UpdatedAt = f.Clock.Now()

	return m
}
