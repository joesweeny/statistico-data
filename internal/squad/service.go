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
			res, err := s.Client.TeamsBySeasonId(id, []string{}, 5)

			if err != nil {
				log.Printf("Error when calling client. Message: %s", err.Error())
				return
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

		_, err := s.Repository.BySeasonAndTeam(seasonId, team.ID)

		if err != ErrNotFound {
			continue
		}

		res, err := s.Client.SquadBySeasonAndTeam(seasonId, team.ID, []string{}, 5)

		if err != nil {
			log.Printf("Error when calling client. Message: %s", err.Error())
			return
		}

		go func(seasonId, teamID int, squad []sportmonks.SquadPlayer) {
			s.persistSquad(seasonId, team.ID, &squad)
			defer waitGroup.Done()
		}(seasonId, team.ID, res.Data)
	}
}

func (s Service) persistSquad(seasonId, teamId int, m *[]sportmonks.SquadPlayer) {
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