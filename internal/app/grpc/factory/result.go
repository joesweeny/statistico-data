package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
)

type ResultFactory struct {
	resultRepo 		app.ResultRepository
	teamRepo        app.TeamRepository
	venueRepo       app.VenueRepository
	logger 			*logrus.Logger
}

func (r ResultFactory) BuildResult(f *app.Fixture) (*proto.Result, error) {
	x, err := r.resultRepo.ByFixtureID(f.ID)

	if err != nil {
		return nil, r.returnLoggedError(f.ID, err)
	}

	home, err := r.teamRepo.ByID(f.HomeTeamID)

	if err != nil {
		return nil, r.returnLoggedError(f.ID, err)
	}

	away, err := r.teamRepo.ByID(f.AwayTeamID)

	if err != nil {
		return nil, r.returnLoggedError(f.ID, err)

	}

	p := proto.Result{
		Id:        int64(x.FixtureID),
		DateTime:  f.Date.Unix(),
		MatchData: toMatchData(home, away, x),
	}

	if f.VenueID != nil {
		v, err := r.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = venueToProto(v)
		}
	}

	return &p, nil
}

func (r ResultFactory) returnLoggedError(id uint64, err error) error {
	r.logger.Warnf("error when hydrating proto result: fixture %d. error %s", id, err.Error())
	return err
}

func NewResultFactory(r app.ResultRepository, t app.TeamRepository, v app.VenueRepository, log *logrus.Logger) *ResultFactory {
	return &ResultFactory{resultRepo: r, teamRepo: t, venueRepo: v, logger: log}
}
