package sportmonks

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
)

type CountryRequester struct {
	client *sportmonks.Client
	clock clockwork.Clock
}

func (c CountryRequester) Countries(ch chan<- *app.Country) error {
	res, err := c.callClient(1, ch)

	if err != nil {
		return err
	}

	total := res.Meta.Pagination.Total

	if total > 1 {
		for i := 2; i <= total; i++ {
			_, err := c.callClient(i, ch)

			if err != nil {
				return err
			}
		}
	}

	close(ch)

	return nil
}

func (c CountryRequester) callClient(page int, ch chan<- *app.Country) (*sportmonks.CountriesResponse, error) {
	res, err := c.client.Countries(page, []string{}, 5)

	if err != nil {
		return &sportmonks.CountriesResponse{}, err
	}

	for _, country := range res.Data {
		ch <- c.hydrateCountry(&country)
	}

	return res, nil
}

func (c CountryRequester) hydrateCountry(s *sportmonks.Country) *app.Country {
	return &app.Country{
		ID:        s.ID,
		Name:      s.Name,
		Continent: s.Extra.Continent,
		ISO:       s.Extra.ISO,
	}
}

func NewCountryRequester(client *sportmonks.Client, clock clockwork.Clock) *CountryRequester {
	return &CountryRequester{
		client: client,
		clock:  clock,
	}
}