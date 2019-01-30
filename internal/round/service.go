package round

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
	q := []string{"rounds"}

	for _, id := range ids {
		res, err := s.Client.SeasonById(id, q)

		if err != nil {
			return err
		}

		for _, round := range res.Data.Rounds.Data {
			// Push method into Go routine
			s.persistRound(&round)
		}
	}

	return nil
}

func (s Service) persistRound(m *sportmonks.Round) {
	round, err := s.GetById(m.ID)

	if err != nil && (model.Round{}) == *round {
		created, err := s.createRound(m)

		if err != nil {
			log.Printf("Error occurred when creating struct: %s", err.Error())
		}

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when inserting struct %+v", created)
		}

		return
	}

	updated, err := s.updateRound(m, round)

	if err != nil {
		log.Printf("Error occurred when updating struct: %s", err.Error())
		return
	}

	if err := s.Update(updated); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}
