package season

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"log"
)

type Processor struct {
	Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

const season = "season"

func (s Processor) Process(command string, done chan bool) {
	if command != season {
		s.Logger.Fatalf("Command %s is not supported", command)
		return
	}

	res, err := s.Client.Seasons(1, []string{}, 5)

	if err != nil {
		s.Logger.Fatalf("Error when calling client '%s", err.Error())
		return
	}

	seasons := make(chan sportmonks.Season, res.Meta.Pagination.Total)

	go s.parseSeasons(seasons, res.Meta)
	go s.persistSeasons(seasons, done)

	return
}

func (s Processor) parseSeasons(ch chan<- sportmonks.Season, meta sportmonks.Meta) {
	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Seasons(i, []string{}, 5)

		if err != nil {
			s.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		for _, season := range res.Data {
			ch <- season
		}
	}

	close(ch)
}

func (s Processor) persistSeasons(ch <-chan sportmonks.Season, done chan bool) {
	for x := range ch {
		s.persist(&x)
	}

	done <- true
}

func (s Processor) persist(m *sportmonks.Season) {
	season, err := s.Id(m.ID)

	if err != nil && (model.Season{} == *season) {
		created := s.createSeason(m)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Season struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateSeason(m, season)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s'occurred when updating Season struct: %+v\n,", err.Error(), updated)
	}

	return
}
