package season

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"log"
	"fmt"
	"sync"
)

var waitGroup sync.WaitGroup

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

	for i := res.Meta.Pagination.CurrentPage; i <= res.Meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Seasons(i, []string{})

		if err != nil {
			return err
		}

		for _, season := range res.Data {
			waitGroup.Add(1)

			go func(season sportmonks.Season) {
				s.persistSeason(&season)
				defer waitGroup.Done()
			}(season)
		}
	}

	waitGroup.Wait()

	return nil
}

func (s Service) persistSeason(m *sportmonks.Season) {
	fmt.Printf("Season ID in persist method: %d\n", m.ID)
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
