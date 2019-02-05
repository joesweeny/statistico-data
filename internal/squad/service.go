package squad

import (
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/sportmonks-go-client"
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

	for _, id := range ids {
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.TeamsBySeasonId(id, []string{"squad"}, 5)

			if err != nil {
				log.Printf("Error when calling client. Message: %s", err.Error())
			}

			s.handleTeams(id, res.Data)

			defer waitGroup.Done()
		}(id)
	}

	waitGroup.Wait()

	return nil
}

func (s Service) handleTeams(seasonId int, t []sportmonks.Team) {
	for _, team := range t {
		waitGroup.Add(1)

		go func(seasonId int, team sportmonks.Team) {
			s.persistSquad(seasonId, team.ID, &team.Squad)
			defer waitGroup.Done()
		}(seasonId, team)
	}
}

func (s Service) persistSquad(seasonId, teamId int, m *sportmonks.Squad) {
	squad, err := s.BySeasonAndTeam(seasonId, teamId)

	if err == ErrNotFound {
		created := s.createSquad(seasonId, teamId, m)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Squad struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateSquad(m, squad)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Squad struct: %+v\n,", err.Error(), updated)
	}

	return
}