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

const team = "team"
const teamCurrentSeason = "team:current-season"

var waitGroup sync.WaitGroup

func (s Service) Process(command string, done chan bool) {
	switch command {
	case team:
		go s.allSeasons(done)
	case teamCurrentSeason:
		go s.currentSeason(done)
	default:
		s.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (s Service) allSeasons(done chan bool) {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	teams := make(chan sportmonks.Team, len(ids))

	go s.parseTeamsSync(teams, ids)
	go s.persistTeams(teams, done)
}

func (s Service) currentSeason(done chan bool) {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	s.parseTeamsAsync(ids)

	waitGroup.Wait()

	done <- true
}

func (s Service) parseTeamsSync(ch chan<- sportmonks.Team, ids []int) {
	for _, id := range ids {
		res, err := s.Client.TeamsBySeasonId(id, []string{}, 5)

		if err != nil {
			log.Printf("Error when calling client. Message: %s", err.Error())
		}

		for _, team := range res.Data {
			ch <- team
		}
	}

	close(ch)
}

func (s Service) parseTeamsAsync(ids []int) {
	for _, id := range ids {
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.TeamsBySeasonId(id, []string{}, 5)

			if err != nil {
				log.Printf("Error when calling client. Message: %s", err.Error())
			}

			for _, team := range res.Data {
				waitGroup.Add(1)
				go func(team sportmonks.Team) {
					s.persistTeam(&team)
					defer waitGroup.Done()
				}(team)
			}

			defer waitGroup.Done()
		}(id)
	}
}

func (s Service) persistTeams(ch <-chan sportmonks.Team, done chan bool) {
	for team := range ch {
		s.persistTeam(&team)
	}

	done <- true
}

func (s Service) persistTeam(t *sportmonks.Team) {
	team, err := s.GetById(t.ID)

	if (err == ErrNotFound && model.Team{} == *team) {
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
