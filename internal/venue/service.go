package venue

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/statshub/internal/season"
	"log"
)

type Service struct {
	Repository
	SeasonRepo season.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

func (s Service) Process() error {
	ids, err := s.SeasonRepo.Ids()

	if err != nil {
		return err
	}

	return s.callClient(ids)
}

func (s Service) callClient(ids []int) error {
	for _, id := range ids {
		res, err := s.Client.BySeasonId(id)

		if err != nil {
			return err
		}

		for _, venue := range res.Data {
			// Push method into Go routine
			s.persistVenue(&venue)
		}
	}

	return nil
}

func (s Service) persistVenue(v *sportmonks.Venue) {
	venue, err := s.GetById(v.ID)

	if err != nil && (model.Venue{}) == *venue {
		created := s.createVenue(v)

		if err := s.Insert(created); err != nil {
			log.Printf("Error occurred when creating struct %+v", created)
		}

		return
	}

	updated := s.updateVenue(v, venue)

	if err := s.Update(venue); err != nil {
		log.Printf("Error occurred when updating struct: %+v, error %+v", updated, err)
	}

	return
}
