package bootstrap

import (
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/country"
	"github.com/joesweeny/statshub/internal/fixture"
	"github.com/joesweeny/statshub/internal/round"
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/statshub/internal/venue"
)

type Service interface {
	Process() error
}

func (b Bootstrap) CompetitionService() competition.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	return competition.Service{
		Repository: &competition.PostgresCompetitionRepository{Connection: conn},
		Factory:    competition.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}
}

func (b Bootstrap) CountryService() country.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	c := country.Service{
		Repository: &country.PostgresCountryRepository{Connection: conn},
		Factory:    country.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}

	return c
}

func (b Bootstrap) FixtureService() fixture.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	c := fixture.Service{
		Repository: &fixture.PostgresFixtureRepository{Connection: conn},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: conn},
		Factory:    fixture.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}

	return c
}

func (b Bootstrap) RoundService() round.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	s := round.Service{
		Repository: &round.PostgresRoundRepository{Connection: conn},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: conn},
		Factory:    round.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}

	return s
}

func (b Bootstrap) SeasonService() season.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	c := season.Service{
		Repository: &season.PostgresSeasonRepository{Connection: conn},
		Factory:    season.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}

	return c
}

func (b Bootstrap) VenueService() venue.Service {
	conn := b.databaseConnection()
	client, err := b.sportmonksClient()

	if err != nil {
		panic(err.Error())
	}

	s := venue.Service{
		Repository: &venue.PostgresVenueRepository{Connection: conn},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: conn},
		Factory:    venue.Factory{Clock: clock()},
		Client:     client,
		Logger:     logger(),
	}

	return s
}
