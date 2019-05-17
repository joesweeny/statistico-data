package fixture

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/season"
	"log"
	"sync"
)

var waitGroup sync.WaitGroup

type Processor struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

const fixture = "fixture"
const fixtureCurrentSeason = "fixture:current-season"

func (s Processor) Process(command string, option string, done chan bool) {
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

func (s Processor) allSeasons(done chan bool) {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go s.callClient(ids, done)
}

func (s Processor) currentSeason(done chan bool) {
	ids, err := s.SeasonRepo.CurrentSeasonIds()

	if err != nil {
		s.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go s.callClient(ids, done)
}

func (s Processor) callClient(ids []int, done chan bool) {
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

func (s Processor) handleFixtures(f []sportmonks.Fixture) {
	for _, fixture := range f {
		waitGroup.Add(1)

		go func(fixture sportmonks.Fixture) {
			s.persistFixture(&fixture)
			defer waitGroup.Done()
		}(fixture)
	}
}

func (s Processor) persistFixture(m *sportmonks.Fixture) {
	fixture, err := s.ById(uint64(m.ID))

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
