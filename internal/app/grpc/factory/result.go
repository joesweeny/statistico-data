package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"time"
)

type ResultFactory struct {
	resultRepo app.ResultRepository
	roundRepo app.RoundRepository
	seasonRepo app.SeasonRepository
	teamRepo   app.TeamRepository
	venueRepo  app.VenueRepository
	logger     *logrus.Logger
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

	season, err := r.seasonRepo.ByID(f.SeasonID)

	if err != nil {
		return nil, r.returnLoggedError(f.ID, err)
	}

	p := proto.Result{
		Id:        x.FixtureID,
		HomeTeam:  TeamToProto(home),
		AwayTeam:  TeamToProto(away),
		Season:    SeasonToProto(season),
		DateTime:  &proto.Date{
			Utc: f.Date.Unix(),
			Rfc: f.Date.Format(time.RFC3339),
		},
		Stats: toMatchStats(x),
	}

	if f.VenueID != nil {
		v, err := r.venueRepo.GetById(*f.VenueID)

		if err == nil {
			p.Venue = venueToProto(v)
		}
	}

	if f.RoundID != nil {
		r, err := r.roundRepo.ByID(*f.RoundID)

		if err == nil {
			p.Round = roundToProto(r)
		}
	}

	return &p, nil
}

func (r ResultFactory) returnLoggedError(id uint64, err error) error {
	r.logger.Warnf("error hydrating proto result: fixture %d. error %s", id, err.Error())
	return err
}

func NewResultFactory(
	r app.ResultRepository,
	o app.RoundRepository,
	s app.SeasonRepository,
	t app.TeamRepository,
	v app.VenueRepository,
	log *logrus.Logger,
) *ResultFactory {
	return &ResultFactory{resultRepo: r, roundRepo: o, seasonRepo: s, teamRepo: t, venueRepo: v, logger: log}
}
