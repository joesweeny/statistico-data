package season

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"github.com/joesweeny/statshub/internal/model"
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

	for i := res.Meta.Pagination.CurrentPage; i <= res.Meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Seasons(i, []string{})

		if err != nil {
			return err
		}

		for _, season := range res.Data {
			// Push method into a Go routine
			s.persistSeason(&season)
		}
	}

	return nil
}

func (s Service) persistSeason(m *sportmonks.Season) {
	season, err := s.GetById(m.ID)

	if err != nil && (model.Season{} == *season) {
		created := s.createSeason(m)

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateSeason(m, season)

	if err := s.Update(updated); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}