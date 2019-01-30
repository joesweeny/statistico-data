package fixture

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/statshub/internal/season"
	"log"
)

type Service struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

func (s Service) Process() error {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		return err
	}

	return s.callClient(ids)
}

func (s Service) CurrentSeason() error {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		return err
	}

	return s.callClient(ids)
}

func (s Service) callClient(ids []int) error {
	q := []string{"fixtures"}

	for _, id := range ids {
		res, err := s.Client.SeasonById(id, q)

		if err != nil {
			return err
		}

		for _, fixture := range res.Data.Fixtures.Data {
			// Push method into Go routine
			s.persistFixture(&fixture)
		}
	}

	return nil
}

func (s Service) persistFixture(m *sportmonks.Fixture) {
	fixture, err := s.GetById(m.ID)

	if err != nil && (model.Fixture{}) == *fixture {
		created := s.createFixture(m)

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateFixture(m, fixture)

	if err := s.Update(updated); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}
