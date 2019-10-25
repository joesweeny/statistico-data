package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const country = "country"

// CountryProcessor is used to process data from an external data source to this applications
// chosen data store
type CountryProcessor struct {
	repository app.CountryRepository
	requester  app.CountryRequester
	logger     *logrus.Logger
}

// Process fetches data from external an external data source using the CountryRequester
// before persisting to the storage engine using the CountryRepository
func (p CountryProcessor) Process(command string, option string, done chan bool) {
	if command != country {
		p.logger.Fatalf("Command %s is not supported", command)
	}

	ch := p.requester.Countries()

	go p.persistCountries(ch, done)
}

func (p CountryProcessor) persistCountries(ch <-chan *app.Country, done chan bool) {
	for country := range ch {
		p.persist(country)
	}

	done <- true
}

func (p CountryProcessor) persist(c *app.Country) {
	_, err := p.repository.GetById(c.ID)

	if err != nil {
		if err := p.repository.Insert(c); err != nil {
			p.logger.Warningf("Error '%s' occurred when inserting country struct: %+v\n,", err.Error(), *c)
		}

		return
	}

	if err := p.repository.Update(c); err != nil {
		p.logger.Warningf("Error '%s' occurred when updating country struct: %+v\n,", err.Error(), *c)
	}

	return
}

func NewCountryProcessor(r app.CountryRepository, s app.CountryRequester, log *logrus.Logger) *CountryProcessor {
	return &CountryProcessor{repository: r, requester: s, logger: log}
}
