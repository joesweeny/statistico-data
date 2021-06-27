package rest

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"time"
)

type FixtureFactory struct {
	roundRepo app.RoundRepository
	teamRepo  app.TeamRepository
	venueRepo app.VenueRepository
}

func (b FixtureFactory) BuildFixture(f *app.Fixture) (*Fixture, error) {
	home, err := b.teamRepo.ByID(f.HomeTeamID)

	if err != nil {
		return nil, err
	}

	away, err := b.teamRepo.ByID(f.AwayTeamID)

	if err != nil {
		return nil, err
	}

	p := Fixture{
		ID:       f.ID,
		HomeTeam: convertAppTeam(home),
		AwayTeam: convertAppTeam(away),
		Date: Date{
			UTC: uint64(f.Date.Unix()),
			RFC: f.Date.Format(time.RFC3339),
		},
	}

	if f.VenueID != nil {
		v, err := b.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = convertAppVenue(v)
		}
	}

	if f.RoundID != nil {
		r, err := b.roundRepo.ByID(*f.RoundID)

		if err == nil {
			p.Round = convertAppRound(r)
		}
	}

	return &p, nil
}

func NewFixtureFactory(r app.RoundRepository, t app.TeamRepository, v app.VenueRepository) *FixtureFactory {
	return &FixtureFactory{roundRepo: r, teamRepo: t, venueRepo: v}
}
