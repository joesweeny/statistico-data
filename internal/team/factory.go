package team

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createTeam(s *sportmonks.Team) *model.Team {
	return &model.Team{
		ID:           s.ID,
		Name:         s.Name,
		ShortCode:    &s.ShortCode,
		CountryID:    &s.CountryID,
		VenueID:      s.VenueID,
		NationalTeam: s.NationalTeam,
		Founded:      &s.Founded,
		Logo:         s.LogoPath,
		CreatedAt:    f.Clock.Now(),
		UpdatedAt:    f.Clock.Now(),
	}
}

func (f Factory) updateTeam(s *sportmonks.Team, m *model.Team) *model.Team {
	m.Name = s.Name
	m.ShortCode = &s.ShortCode
	m.VenueID = s.VenueID
	m.Logo = s.LogoPath
	m.UpdatedAt = f.Clock.Now()

	return m
}
