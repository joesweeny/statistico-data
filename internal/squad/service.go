package squad

import (
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/sportmonks-go-client"
	"log"
)

const callLimit = 2000

type Service struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

var counter int

func (s Service) Process() error {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		return err
	}

	done := make(chan bool)

	go s.handleSeasons(ids, done, &counter)

	<-done

	return nil
}

func (s Service) handleSeasons(ids []int, done chan bool, c *int) {
	teams := make(chan sportmonks.Team, callLimit)

	for _, id := range ids {
		if *c >= callLimit {
			s.Logger.Printf("Api call limited reached %d calls\n", *c)
			done <- true
		}

		res, err := s.Client.TeamsBySeasonId(id, []string{}, 5)

		if err != nil {
			s.Logger.Printf("Error when calling client. Message: %s\n", err.Error())
			return
		}

		for _, t := range res.Data {
			teams <- t
		}

		go s.handleTeams(id, teams, c, done)
	}
}

func (s Service) handleTeams(seasonId int, teams chan sportmonks.Team, c *int, done chan bool) {
	for t := range teams {
		if *c >= callLimit {
			s.Logger.Printf("Api call limited reached %d calls\n", *c)
			done <- true
		}

		_, err := s.BySeasonAndTeam(seasonId, t.ID)

		if err != ErrNotFound {
			continue
		}

		res, err := s.Client.SquadBySeasonAndTeam(seasonId, t.ID, []string{}, 5)

		*c++

		if err != nil {
			s.Logger.Printf("Error when calling client. Message: %s", err.Error())
			return
		}

		go s.persistSquad(seasonId, t.ID, &res.Data)
	}
}

func (s Service) persistSquad(seasonId, teamId int, m *[]sportmonks.SquadPlayer) {
	_, err := s.BySeasonAndTeam(seasonId, teamId)

	if err == ErrNotFound {
		created := s.createSquad(seasonId, teamId, m)

		if err := s.Insert(created); err != nil {
			s.Logger.Printf("Error '%s' occurred when inserting Squad struct: %+v\n,", err.Error(), created)
		}

		return
	}

	return
}
