package bootstrap

import (
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/country"
)

type Service interface {
	Process() error
}

func (b Bootstrap) GetCompetitionService() competition.Service {
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

func (b Bootstrap) GetCountryService() country.Service {
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
