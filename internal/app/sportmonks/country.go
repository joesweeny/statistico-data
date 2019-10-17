package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type CountryRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (c CountryRequester) Countries() <-chan *app.Country {
	_, meta, err := c.client.Countries(context.Background(),1, []string{})

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s' when making country request", err.Error())
		return nil
	}

	ch := make(chan *app.Country, meta.Pagination.Total)

	go c.parseCountries(meta.Pagination.TotalPages, ch)

	return ch
}

func (c CountryRequester) parseCountries(pages int, ch chan<- *app.Country) {
	defer close(ch)

	for i := 1; i <= pages; i++ {
		c.callClient(i, ch)
	}
}

func (c CountryRequester) callClient(page int, ch chan<- *app.Country) {
	res, _, err := c.client.Countries(context.Background(), page, []string{})

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s' when making country request", err.Error())
		return
	}

	for _, country := range res {
		ch <- transform(&country)
	}
}

func transform(s *spClient.Country) *app.Country {
	return &app.Country{
		ID:        int64(s.ID),
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
	}
}

func NewCountryRequester(client *spClient.HTTPClient, log *logrus.Logger) *CountryRequester {
	return &CountryRequester{client: client, logger: log}
}
