package competition

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
	res, err := s.Client.Leagues(1, []string{})

	if err != nil {
		return err
	}

	comps := make(chan sportmonks.League, res.Meta.Pagination.Total)
	done := make(chan bool)

	go s.parseLeagues(comps, res.Meta)
	go s.persistCompetitions(comps, done)

	<-done

	return nil
}

func (s Service) parseLeagues(ch chan<- sportmonks.League, meta sportmonks.Meta) {
	for i := meta.Pagination.CurrentPage; i <= meta.Pagination.TotalPages; i++ {
		res, err := s.Client.Leagues(i, []string{})

		if err != nil {
			log.Printf("Error when calling client '%s", err.Error())
		}

		for _, comp := range res.Data {
			ch <- comp
		}
	}

	close(ch)
}

func (s Service) persistCompetitions(ch <-chan sportmonks.League, done chan bool) {
	for x := range ch {
		s.persist(&x)
	}

	done <- true
}

func (s Service) persist(l *sportmonks.League) {
	comp, err := s.GetById(l.ID)

	if err != nil && (model.Competition{}) == *comp {
		created := s.createCompetition(l)

		if err := s.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Competition struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := s.updateCompetition(l, comp)

	if err := s.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Competition struct: %+v\n,", err.Error(), updated)
	}

	return
}
