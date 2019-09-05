package process

import (
	spClient "github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/statistico"
	"github.com/statistico/statistico-data/internal/statistico/sportmonks"
	"log"
)

// CountryProcessor is used to retrieve data from an external data source and
// persist to this applications chosen data store
type CountryProcessor struct {
	repository statistico.CountryRepository
	factory sportmonks.CountryFactory
	client *spClient.Client
	logger *log.Logger
}

const country = "country"

func (p CountryProcessor) Process(command string, option string, done chan bool) {
	if command != country {
		p.logger.Fatalf("Command %s is not supported", command)
		return
	}

	res, err := p.client.Countries(1, []string{}, 5)

	if err != nil {
		p.logger.Fatalf("Error when calling client '%s", err.Error())
		return
	}

	countries := make(chan spClient.Country, res.Meta.Pagination.Total)

	go p.parseCountries(countries, res.Meta)
	go p.persistCountries(countries, done)
}

func (p CountryProcessor) parseCountries(ch chan<- spClient.Country, meta spClient.Meta) {
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

	close(ch)
}

func (p CountryProcessor) persistCountries(ch <-chan spClient.Country, done chan bool) {
	for x := range ch {
		p.persist(&x)
	}

	done <- true
}

func (p CountryProcessor) persist(c *spClient.Country) {
	country, err := p.repository.GetById(c.ID)

	if err != nil {
		created := p.factory.CreateCountry(c)

		if err := p.repository.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Country struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := p.factory.UpdateCountry(c, country)

	if err := p.repository.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Competition struct: %+v\n,", err.Error(), updated)
	}

	return
}

