package bootstrap

import (
	"github.com/joesweeny/statshub/internal/country"
)

func (b Bootstrap) GetCountryService() (country.Service) {
	conn := databaseConnection()
	client, err := sportmonksClient()

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
