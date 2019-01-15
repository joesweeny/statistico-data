package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
)

type Service struct {
	client     sportmonks.Client
	handler	   Handler
	logger     log.Logger
}

func (s Service) HandleCountries() error {
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
			// Handle error, leaving blank for now - push method into a Go routine
			// Log error out
			s.handler.Handle(country)
		}

		i++
	}

	return nil
}
