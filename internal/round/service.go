package round

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/statshub/internal/season"
	"log"
	"sync"
)

var waitGroup sync.WaitGroup

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
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.SeasonById(id, q)

			if err != nil {
				log.Printf("Error when calling client '%s", err.Error())
			}

			s.handleRounds(res.Data.Rounds.Data)

			defer waitGroup.Done()
		}(id)
	}

	waitGroup.Wait()

	return nil
}

func (s Service) handleRounds(r []sportmonks.Round) {
	for _, round := range r {
		waitGroup.Add(1)

		go func(round sportmonks.Round) {
			s.persistRound(&round)
			defer waitGroup.Done()
		}(round)
	}
}

func (s Service) persistRound(m *sportmonks.Round) {
	round, err := s.GetById(m.ID)

	if err != nil && (model.Round{}) == *round {
		created, err := s.createRound(m)

		if err != nil {
			log.Printf("Error '%s' occurred when creating Round struct: %+v\n,", err.Error(), created)
			return
		}

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Round struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated, err := s.updateRound(m, round)

	if err != nil {
		log.Printf("Error '%s' occurred when updating Round struct: %+v\n,", err.Error(), updated)
		return
	}

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Round struct: %+v\n,", err.Error(), updated)
	}

	return
}