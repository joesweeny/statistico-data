package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"time"
)

type FixtureFactory struct {
	roundRepo       app.RoundRepository
	teamRepo        app.TeamRepository
	venueRepo       app.VenueRepository
	logger 			*logrus.Logger
}

func (b FixtureFactory) BuildFixture(f *app.Fixture) (*proto.Fixture, error) {
	home, err := b.teamRepo.ByID(f.HomeTeamID)

	if err != nil {
		return nil, b.returnLoggedError(f.ID, err)
	}

	away, err := b.teamRepo.ByID(f.AwayTeamID)

	if err != nil {
		return nil, b.returnLoggedError(f.ID, err)
	}

	p := proto.Fixture{
		Id:       int64(f.ID),
		HomeTeam: TeamToProto(home),
		AwayTeam: TeamToProto(away),
		DateTime:    &proto.Date{
			Utc: f.Date.Unix(),
			Rfc: f.Date.Format(time.RFC3339),
		},
	}

	if f.VenueID != nil {
		v, err := b.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = VenueToProto(v)
		}
	}

	if f.RoundID != nil {
		r, err := b.roundRepo.ByID(*f.RoundID)

		if err == nil {
			p.Round = RoundToProto(r)
		}
	}

	return &p, nil
}

func (b FixtureFactory) returnLoggedError(id uint64, err error) error {
	b.logger.Warnf("error when hydrating proto fixture: fixture %d. error %s", id, err.Error())
	return err
}

func NewFixtureFactory(r app.RoundRepository, t app.TeamRepository, v app.VenueRepository, log *logrus.Logger) *FixtureFactory {
	return &FixtureFactory{roundRepo: r, teamRepo: t, venueRepo: v, logger: log}
}
