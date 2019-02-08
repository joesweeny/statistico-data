package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"log"
)

type Processor struct {
	Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

const country = "country"

func (s Processor) Process(command string, done chan bool) {
	if command != country {
		s.Logger.Fatalf("Command %s is not supported", command)
		return
	}

	res, err := s.Client.Countries(1, []string{}, 5)

	if err != nil {
		s.Logger.Fatalf("Error when calling client '%s", err.Error())
		return
	}

	countries := make(chan sportmonks.Country, res.Meta.Pagination.Total)

	go s.parseCountries(countries, res.Meta)
	go s.persistCountries(countries, done)
}

func (s Processor) parseCountries(ch chan<- sportmonks.Country, meta sportmonks.Meta) {
	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Countries(i, []string{}, 5)

		if err != nil {
			s.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		for _, country := range res.Data {
			ch <- country
		}
	}

	close(ch)
}

func (s Processor) persistCountries(ch <-chan sportmonks.Country, done chan bool) {
	for x := range ch {
		s.persist(&x)
	}

	done <- true
}

func (s Processor) persist(c *sportmonks.Country) {
	country, err := s.GetById(c.ID)

	if err != nil && (model.Country{}) == *country {
		created := s.createCountry(c)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Country struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateCountry(c, country)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Competition struct: %+v\n,", err.Error(), updated)
	}

	return
}
