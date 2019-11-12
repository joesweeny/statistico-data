package fetch

import (
	"fmt"
	"github.com/statistico/statistico-data/internal/app"
)

type EntityFetcher struct {
	competitionRepo app.CompetitionRepository
	resultRepo      app.ResultRepository
	roundRepo       app.RoundRepository
	seasonRepo      app.SeasonRepository
	teamRepo        app.TeamRepository
	venueRepo       app.VenueRepository
}

func (e EntityFetcher) ResultByID(id uint64) (*app.Result, error) {
	r, err := e.resultRepo.ByFixtureID(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving result %d in EntityFetcher", id)
		return nil, e
	}

	return r, nil
}

func (e EntityFetcher) SeasonByID(id uint64) (*app.Season, error) {
	s, err := e.seasonRepo.ByID(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving season %d in EntityFetcher", id)
		return nil, e
	}

	return s, nil
}

func (e EntityFetcher) CompetitionByID(id uint64) (*app.Competition, error) {
	c, err := e.competitionRepo.ByID(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving competition %d in EntityFetcher", id)
		return nil, e
	}

	return c, nil
}

func (e EntityFetcher) TeamByID(id uint64) (*app.Team, error) {
	t, err := e.teamRepo.ByID(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving team %d in EntityFetcher", id)
		return nil, e
	}

	return t, nil
}

func (e EntityFetcher) RoundByID(id uint64) (*app.Round, error) {
	r, err := e.roundRepo.ByID(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving round %d in EntityFetcher", id)
		return nil, e
	}

	return r, nil
}

func (e EntityFetcher) VenueByID(id uint64) (*app.Venue, error) {
	v, err := e.venueRepo.GetById(id)

	if err != nil {
		e := fmt.Errorf("error when retrieving venue %d in EntityFetcher", id)
		return nil, e
	}

	return v, nil
}
