package season

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"log"
)

type Service struct {
	Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

func (s Service) Process() error {
	res, err := s.Client.Seasons(1, []string{})

	if err != nil {
		return err
	}

	seasons := make(chan sportmonks.Season, res.Meta.Pagination.Total)
	done := make(chan bool)

	go s.parseSeasons(seasons, res.Meta)
	go s.persistSeasons(seasons, done)

	<- done

	return nil
}

func (s Service) parseSeasons(ch chan<- sportmonks.Season, meta sportmonks.Meta) {
	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Seasons(i, []string{})

		if err != nil {
			log.Printf("Error when calling client '%s", err.Error())
		}

		for _, season := range res.Data {
			ch <- season
		}
	}

	close(ch)
}

func (s Service) persistSeasons(ch <-chan sportmonks.Season, done chan bool) {
	for x := range ch {
		s.persist(&x)
	}

	done <- true
}

func (s Service) persist(m *sportmonks.Season) {
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
