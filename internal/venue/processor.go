package venue

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

const venue = "venue"
const venueCurrentSeason = "venue:current-season"

func (s Processor) Process(command string, done chan bool) {
	switch command {
	case venue:
		go s.allSeasons(done)
	case venueCurrentSeason:
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
	for _, id := range ids {
		waitGroup.Add(1)

		go func(i int) {
			res, err := s.Client.VenuesBySeasonId(i, 5)

			if err != nil {
				log.Printf("Error when calling client '%s", err.Error())
			}

			s.handleVenues(res.Data)

			defer waitGroup.Done()
		}(id)
	}

	waitGroup.Wait()

	done <- true
}

func (s Processor) handleVenues(v []sportmonks.Venue) {
	for _, venue := range v {
		waitGroup.Add(1)

		go func(ven sportmonks.Venue) {
			s.persistVenue(&ven)
			defer waitGroup.Done()
		}(venue)
	}
}

func (s Processor) persistVenue(v *sportmonks.Venue) {
	venue, err := s.GetById(v.ID)

	if err != nil && (model.Venue{}) == *venue {
		created := s.createVenue(v)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Venue struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateVenue(v, venue)

	if err := s.Update(venue); err != nil {
		log.Printf("Error '%s' occurred when updating Venue struct: %+v\n,", err.Error(), updated)
	}

	return
}
