package country

import (
	"github.com/satori/go.uuid"
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type factory struct {
	clock clockwork.Clock
}

func (f factory) createCountry(s sportmonks.Country, id uuid.UUID) model.Country {
	return model.Country{
		ID: 		id,
		ExternalID: s.ID,
		Name: 		s.Name,
		Continent: 	s.Extra.Continent,
		ISO: 		s.Extra.ISO,
		CreatedAt: 	f.clock.Now(),
		UpdatedAt: 	f.clock.Now(),
	}
}

func (f factory) updateCountry(s sportmonks.Country, m model.Country) model.Country {
	m.ExternalID = s.ID
	m.Name = s.Name
	m.Continent = s.Extra.Continent
	m.ISO = s.Extra.ISO
	m.UpdatedAt = f.clock.Now()

	return m
}
