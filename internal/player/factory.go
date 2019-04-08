package player

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createPlayer(s *sportmonks.Player) *model.Player {
	return &model.Player{
		ID:          s.PlayerID,
		CountryId:   s.CountryID,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		BirthPlace:  &s.Birthplace,
		DateOfBirth: &s.BirthDate,
		PositionID:  s.PositionID,
		Image:       &s.ImagePath,
		CreatedAt:   f.Clock.Now(),
		UpdatedAt:   f.Clock.Now(),
	}
}
