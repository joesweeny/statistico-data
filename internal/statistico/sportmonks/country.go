package sportmonks

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/statistico"
	"log"
)

type CountryDataService struct {
	client *sportmonks.Client
	logger *log.Logger
	clock clockwork.Clock
}

func (c CountryDataService) Countries(ch chan<- *statistico.Country) {
	defer close(ch)

	res, err := c.client.Countries(1, []string{}, 5)

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s", err.Error())
		return
	}

	for _, country := range res.Data {
		ch <- c.createCountry(&country)
	}

	total := res.Meta.Pagination.Total



	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := p.client.Countries(i, []string{}, 5)

		if err != nil {
			p.logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		for _, country := range res.Data {
			ch <- country
		}
	}
}

type CountryFactory struct {
	Clock clockwork.Clock
}

func (c CountryDataService) createCountry(s *sportmonks.Country) *statistico.Country {
	return &statistico.Country{
		ID:        s.ID,
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
	}
}
