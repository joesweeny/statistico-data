package country

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type Handler struct {
	repository 	Repository
	factory Factory
}

func (h Handler) Handle(s sportmonks.Country) error {
	c, err := h.repository.GetByExternalId(s.ID)

	if err != nil && (model.Country{}) == c {
		country := h.factory.create(s)

		if err := h.repository.Insert(country); err != nil {
			return err
		}

		return nil
	}

	country := h.factory.update(s, c)

	if err := h.repository.Update(country); err != nil {
		return err
	}

	return nil
}
