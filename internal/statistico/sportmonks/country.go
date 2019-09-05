package sportmonks

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/statistico"
)

type CountryFactory struct {
	Clock clockwork.Clock
}

func (f CountryFactory) CreateCountry(s *sportmonks.Country) *statistico.Country {
	return &statistico.Country{
		ID:        s.ID,
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f CountryFactory) UpdateCountry(s *sportmonks.Country, m *statistico.Country) *statistico.Country {
	m.Name = s.Name
	m.Continent = s.Extra.Continent
	m.ISO = s.Extra.ISO
	m.UpdatedAt = f.Clock.Now()

	return m
}