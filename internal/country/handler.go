package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type Handler struct {
	Repository Repository
}

func (h Handler) Handle(country sportmonks.Country) error {
	country, err := s.Repository.GetByExternalId(c.ID)

	if err != nil && (model.Country{}) == country {
		s.Repository.Insert(create(c))
		return nil
	}

	s.Repository.Update(update(c, country))

	return nil
}
