package squad

import (
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createSquad(seasonId int, teamId int, s sportmonks.Squad) *model.Squad {
	var x []int

	squad := model.Squad{
		SeasonID: seasonId,
		TeamID:   teamId,
		PlayerIDs: x,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}

	for _, player := range s.Data {
		x = append(x, player.PlayerID)
	}

	return &squad
}

func (f Factory) updateSquad(s *sportmonks.Squad, m *model.Squad) *model.Squad {
	var x []int

	for _, player := range s.Data {
		x = append(x, player.PlayerID)
	}

	m.PlayerIDs = x
	m.UpdatedAt = f.Clock.Now()

	return m
}
