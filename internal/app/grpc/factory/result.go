package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-proto/go"
	"time"
)

type ResultFactory struct {
	resultRepo app.ResultRepository
	roundRepo app.RoundRepository
	seasonRepo app.SeasonRepository
	teamRepo   app.TeamRepository
	teamStatsRepo app.TeamStatsRepository
	venueRepo  app.VenueRepository
	logger     *logrus.Logger
}

func (r ResultFactory) BuildResult(f *app.Fixture) (*statistico.Result, error) {
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

	hs, homeErr := r.teamStatsRepo.ByFixtureAndTeam(f.ID, f.HomeTeamID)
	as, awayErr := r.teamStatsRepo.ByFixtureAndTeam(f.ID, f.AwayTeamID)

	if homeErr != nil || awayErr != nil {
		r.logger.Warnf("error hydrating proto result: fixture %d. error %s: %s", f.ID, homeErr, awayErr)
	}

	date := statistico.Date{
		Utc: f.Date.Unix(),
		Rfc: f.Date.Format(time.RFC3339),
	}

	p := statistico.Result{
		Id:            x.FixtureID,
		HomeTeam:      TeamToProto(home),
		AwayTeam:      TeamToProto(away),
		Season:        SeasonToProto(season),
		DateTime:      &date,
		Stats:         toMatchStats(x),
		HomeTeamStats: TeamStatsToProto(hs),
		AwayTeamStats: TeamStatsToProto(as),
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
	ts app.TeamStatsRepository,
	v app.VenueRepository,
	log *logrus.Logger,
) *ResultFactory {
	return &ResultFactory{
		resultRepo: r,
		roundRepo: o,
		seasonRepo: s,
		teamRepo: t,
		teamStatsRepo: ts,
		venueRepo: v,
		logger: log,
	}
}
