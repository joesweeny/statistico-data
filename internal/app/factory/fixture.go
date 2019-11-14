package factory

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/proto"
	"github.com/statistico/statistico-data/internal/app/rest"
	"time"
)

type FixtureFactory struct {
	roundRepo       app.RoundRepository
	teamRepo        app.TeamRepository
	venueRepo       app.VenueRepository
}

func (b FixtureFactory) BuildProtoFixture(f *app.Fixture) (*proto.Fixture, error) {
	home, err := b.teamRepo.ByID(f.HomeTeamID)

	if err != nil {
		return nil, err
	}

	away, err := b.teamRepo.ByID(f.AwayTeamID)

	if err != nil {
		return nil, err
	}

	p := proto.Fixture{
		Id:          int64(f.ID),
		HomeTeam:    converter.TeamToProto(home),
		AwayTeam:    converter.TeamToProto(away),
		DateTime:    &proto.Date{
			Utc:    f.Date.Unix(),
			Rfc:    f.Date.Format(time.RFC3339),
		},
	}

	if f.VenueID != nil {
		v, err := b.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = converter.VenueToProto(v)
		}
	}

	if f.RoundID != nil {
		r, err := b.roundRepo.ByID(*f.RoundID)

		if err == nil {
			p.Round = converter.RoundToProto(r)
		}
	}

	return &p, nil
}

func (b FixtureFactory) BuildRestFixture(f *app.Fixture) (*rest.Fixture, error) {
	home, err := b.teamRepo.ByID(f.HomeTeamID)

	if err != nil {
		return nil, err
	}

	away, err := b.teamRepo.ByID(f.AwayTeamID)

	if err != nil {
		return nil, err
	}

	p := rest.Fixture{
		ID:          f.ID,
		HomeTeam:    converter.TeamToRest(home),
		AwayTeam:    converter.TeamToRest(away),
		Date:    rest.Date{
			UTC: f.Date.Unix(),
			RFC: f.Date.Format(time.RFC3339),
		},
	}

	if f.VenueID != nil {
		v, err := b.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = converter.VenueToRest(v)
		}
	}

	if f.RoundID != nil {
		r, err := b.roundRepo.ByID(*f.RoundID)

		if err == nil {
			p.Round = converter.RoundToRest(r)
		}
	}

	return &p, nil
}

func NewFixtureFactory(r app.RoundRepository, t app.TeamRepository, v app.VenueRepository) *FixtureFactory {
	return &FixtureFactory{roundRepo: r, teamRepo: t, venueRepo: v}
}
