package country

import (
	"github.com/joesweeny/sportmonks-go-client"
)

type Service struct {
	Repository Repository
	Client     sportmonks.Client
}

func (s Service) HandleCountries() error {
	res, err := s.Client.Countries(1, []string{})

	if err != nil {
		return err
	}

	for i := res.Meta.Pagination.CurrentPage; i <= res.Meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Countries(i, []string{})

		if err != nil {
			return err
		}

		//for _, _ := range res.Data {
		//	// Handle error, leaving blank for now - push method into a Go routine
		//
		//}

		i++
	}

	return nil
}


