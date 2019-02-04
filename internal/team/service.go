package team

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/statshub/internal/season"
	"log"
	"sync"
)

type Service struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

var waitGroup sync.WaitGroup

func (s Service) Process() error {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		return err
	}

	return s.callClient(ids)
}

func (s Service) ProcessCurrentSeason() error {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		return err
	}

	return s.callClient(ids)
}

func (s Service) callClient(ids []int) error {
	for _, id := range ids {
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.TeamsBySeasonId(id, 5)

			if err != nil {
				log.Fatalf("Error when calling client. Message: %s", err.Error())
			}

			s.handleTeams(res.Data)

			defer waitGroup.Done()
		}(id)
	}

	waitGroup.Wait()

	return nil
}

func (s Service) handleTeams(t []sportmonks.Team) {
	for _, team := range t {
		waitGroup.Add(1)

		go func(team sportmonks.Team) {
			s.persistTeam(&team)
			defer waitGroup.Done()
		}(team)
	}
}

func (s Service) persistTeam(t *sportmonks.Team) {
	team, err := s.GetById(t.ID)

	if err != nil && (model.Team{} == *team) {
		created := s.createTeam(t)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Team struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateTeam(t, team)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Team struct: %+v\n,", err.Error(), updated)
	}

	return
}
