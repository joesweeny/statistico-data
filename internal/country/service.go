package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"log"
)

type Service struct {
	Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

func (s Service) Process() error {
	res, err := s.Client.Countries(1, []string{})

	if err != nil {
		return err
	}

	countries := make(chan sportmonks.Country, res.Meta.Pagination.Total)
	done := make(chan bool)

	go s.parseCountries(countries, res.Meta)
	go s.persistCountries(countries, done)

	<- done

	return nil
}

func (s Service) parseCountries(ch chan<- sportmonks.Country, meta sportmonks.Meta) {
	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Countries(i, []string{})

		if err != nil {
			log.Printf("Error when calling client '%s", err.Error())
			continue
		}

		for _, country := range res.Data {
			ch <- country
		}
	}

	close(ch)
}

func (s Service) persistCountries(ch <-chan sportmonks.Country, done chan bool) {
	for x := range ch {
		s.persist(&x)
	}

	done <- true
}

func (s Service) persist(c *sportmonks.Country) {
	country, err := s.GetById(c.ID)

	if err != nil && (model.Country{}) == *country {
		created := s.createCountry(c)

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateCountry(c, country)

	if err := s.Update(updated); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}
