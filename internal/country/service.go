package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
	"time"
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

		for _, country := range res.Data {
			// Handle error, leaving blank for now - push method into a Go routine
			s.persistCountry(country)
		}

		i++
	}

	return nil
}

func (s Service) persistCountry(c sportmonks.Country) error {
	country, err := s.Repository.GetByExternalId(c.ID)

	if err != nil && (model.Country{}) == country {
		s.Repository.Insert(create(c))
		return nil
	}

	s.Repository.Update(update(c, country))

	return nil
}

func create(s sportmonks.Country) model.Country {
	return model.Country{
		ID: 		uuid.Must(uuid.NewV4(), nil),
		ExternalID: s.ID,
		Name: 		s.Name,
		Continent: 	s.Extra.Continent,
		ISO: 		s.Extra.ISO,
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
	}
}

func update(s sportmonks.Country, m model.Country) model.Country {
	m.ExternalID = s.ID
	m.Name = s.Name
	m.Continent = s.Extra.Continent
	m.ISO = s.Extra.ISO

	return m
}