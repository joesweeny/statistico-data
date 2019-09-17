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

func (c CountryRequester) Countries() <-chan *app.Country {
	res, err := c.client.Countries(1, []string{}, 5)

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s when make country request", err.Error())
	}

	ch := make(chan *app.Country, 50)

	go c.parseCountries(res.Meta.Pagination.TotalPages, ch)

	return ch
}

func (c CountryRequester) parseCountries(pages int, ch chan<- *app.Country) {
	defer close(ch)

	for i := 1; i <= pages; i++ {
		c.callClient(i, ch)
	}
}

func (c CountryRequester) callClient(page int, ch chan<- *app.Country) {
	res, err := c.client.Countries(page, []string{}, 5)

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s when make country request", err.Error())
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