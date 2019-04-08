package venue

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createVenue(s *sportmonks.Venue) *model.Venue {
	return &model.Venue{
		ID:        s.ID,
		Name:      s.Name,
		Surface:   &s.Surface,
		Address:   s.Address,
		City:      &s.City,
		Capacity:  &s.Capacity,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}
}

func (f Factory) updateVenue(s *sportmonks.Venue, m *model.Venue) *model.Venue {
	m.Name = s.Name
	m.Surface = &s.Surface
	m.Address = s.Address
	m.City = &s.City
	m.Capacity = &s.Capacity
	m.UpdatedAt = f.Clock.Now()

	return m
}
