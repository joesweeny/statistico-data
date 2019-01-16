package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
)

type Service struct {
	repository
	factory
	client     sportmonks.Client
	logger     *log.Logger
}

func (s Service) Process() error {
	res, err := s.client.Countries(1, []string{})

	if err != nil {
		return err
	}

	for i := res.Meta.Pagination.CurrentPage; i <= res.Meta.Pagination.TotalPages; i++ {
		res, err := s.client.Countries(i, []string{})

		if err != nil {
			return err
		}

		for _, country := range res.Data {
			// Push method into a Go routine
			s.persistCountry(country)
		}

		i++
	}

	return nil
}

func (s Service) persistCountry(c sportmonks.Country) {
	country, err := s.getByExternalId(c.ID)

	if err != nil && (model.Country{}) == country {
		created := s.createCountry(c, uuid.Must(uuid.NewV4(), nil))

		if err := s.insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateCountry(c, country)

	if err := s.update(updated); err != nil {
		log.Printf("Error occurred when updating struct %+v", updated)
	}

	return
}


