package sportmonks

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
	"log"
)

type CountryRequester struct {
	client *sportmonks.Client
	logger     *log.Logger
}

func (c CountryRequester) Countries(ch chan<- *app.Country) {
	res, err := c.client.Countries(1, []string{}, 5)

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s", err.Error())
	}

	for i := 1; i <= res.Meta.Pagination.TotalPages; i++ {
		c.callClient(i, ch)
	}

	close(ch)
}

func (c CountryRequester) callClient(page int, ch chan<- *app.Country) {
	res, err := c.client.Countries(page, []string{}, 5)

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s", err.Error())
	}

	for _, country := range res.Data {
		ch <- transform(&country)
	}
}

func transform(s *sportmonks.Country) *app.Country {
	return &app.Country{
		ID:        s.ID,
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
	}
}

func NewCountryRequester(client *sportmonks.Client, log *log.Logger) *CountryRequester {
	return &CountryRequester{client: client, logger: log}
}