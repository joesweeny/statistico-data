package competition

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"github.com/joesweeny/statshub/internal/model"
)

type Service struct {
	Repository
	Factory
	Client 		*sportmonks.Client
	Logger 		*log.Logger
}

func (s Service) Process() error {
	res, err := s.Client.Leagues(1, []string{})

	if err != nil {
		return err
	}

	for i := res.Meta.Pagination.CurrentPage; i <= res.Meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Leagues(i, []string{})

		if err != nil {
			return err
		}

		for _, comp := range res.Data {
			// Push method into a Go routine
			s.persistCompetition(&comp)
		}
	}

	return nil
}

func (s Service) persistCompetition(l *sportmonks.League) {
	comp, err := s.GetById(l.ID)

	if err != nil && (model.Competition{}) == *comp {
		created := s.createCompetition(l)

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateCompetition(l, comp)

	if err := s.Update(updated); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}
