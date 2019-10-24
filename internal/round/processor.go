package round

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

const round = "round"
const roundCurrentSeason = "round:current-season"

func (s Processor) Process(command string, option string, done chan bool) {
	switch command {
	case round:
		go s.allSeasons(done)
	case roundCurrentSeason:
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

func (s Processor) callClient(ids []int64, done chan bool) {
	q := []string{"rounds"}

	for _, id := range ids {
		waitGroup.Add(1)

		go func(id int) {
			res, err := s.Client.SeasonById(id, q, 5)

			if err != nil {
				log.Printf("Error when calling client '%s", err.Error())
			}

			s.handleRounds(res.Data.Rounds.Data)

			defer waitGroup.Done()
		}(int(id))
	}

	waitGroup.Wait()

	done <- true
}

func (s Processor) handleRounds(r []sportmonks.Round) {
	for _, round := range r {
		waitGroup.Add(1)

		go func(round sportmonks.Round) {
			s.persistRound(&round)
			defer waitGroup.Done()
		}(round)
	}
}

func (s Processor) persistRound(m *sportmonks.Round) {
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
