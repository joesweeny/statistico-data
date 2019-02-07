package fixture

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

const fixture = "fixture"
const fixtureCurrentSeason = "fixture:current-season"

func (s Service) Process(command string, done chan bool) {
	switch command {
	case fixture:
		go s.allSeasons(done)
	case fixtureCurrentSeason:
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

	go s.callClient(ids, done)
}

func (s Service) currentSeason(done chan bool) {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go s.callClient(ids, done)
}

func (s Service) callClient(ids []int, done chan bool) {
	q := []string{"fixtures"}

	for _, id := range ids {
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.SeasonById(id, q, 5)

			if err != nil {
				log.Printf("Error when calling client '%s", err.Error())
			}

			s.handleFixtures(res.Data.Fixtures.Data)

			defer waitGroup.Done()
		}(id)
	}

	waitGroup.Wait()

	done <- true
}

func (s Service) handleFixtures(f []sportmonks.Fixture) {
	for _, fixture := range f {
		waitGroup.Add(1)

		go func(fixture sportmonks.Fixture) {
			s.persistFixture(&fixture)
			defer waitGroup.Done()
		}(fixture)
	}
}

func (s Service) persistFixture(m *sportmonks.Fixture) {
	fixture, err := s.GetById(m.ID)

	if err != nil && (model.Fixture{}) == *fixture {
		created := s.createFixture(m)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Fixture struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateFixture(m, fixture)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Fixture struct: %+v\n,", err.Error(), updated)
	}

	return
}
