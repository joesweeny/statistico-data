package squad

import (
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"fmt"
)

const callLimit = 1500
const squad = "squad"
const squadCurrentSeason = "squad:current-season"

type Service struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

var counter int

func (s Service) Process(command string, done chan bool) {
	if command == squad {
		go s.allSquads(done)
	}

	if command == squadCurrentSeason {
		go s.currentSeason(done)
	}

	s.Logger.Fatalf("Command %s is not supported", command)

	return
}

func (s Service) allSquads(done chan bool) {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	teams := make(chan sportmonks.Team, callLimit)

	go s.handleSeasons(ids, teams, done, &counter)
}

func (s Service) currentSeason(done chan bool) {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	teams := make(chan sportmonks.Team, callLimit)

	go s.handleSeasons(ids, teams, done, &counter)
}

func (s Service) handleSeasons(ids []int, teams chan sportmonks.Team, done chan bool, c *int) {
	for _, id := range ids {
		if *c >= callLimit {
			s.Logger.Printf("Api call limited reached %d calls\n", *c)
			done <- true
		}

		res, err := s.Client.TeamsBySeasonId(id, []string{}, 5)

		if err != nil {
			s.Logger.Printf("Error when calling client. Message: %s\n", err.Error())
			done <- true
		}

		for _, t := range res.Data {
			teams <- t
		}

		fmt.Printf("Parsed all teams for Season %d\n", id)

		go s.handleTeams(id, teams, c, done)
	}

	fmt.Println("Not in Season ID loop anymore")

	close(teams)
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
			done <- true
		}

		go s.persistSquad(seasonId, t.ID, &res.Data)
	}

	done <- true
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
