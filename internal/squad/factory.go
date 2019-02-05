package squad

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createSquad(seasonId int, teamId int, s *[]sportmonks.SquadPlayer) *model.Squad {
	squad := model.Squad{
		SeasonID:  seasonId,
		TeamID:    teamId,
		PlayerIDs: []int{},
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}

	for _, player := range *s {
		squad.PlayerIDs = append(squad.PlayerIDs, player.PlayerID)
	}

	return &squad
}

func (f Factory) updateSquad(s *[]sportmonks.SquadPlayer, m *model.Squad) *model.Squad {
	var x []int

	for _, player := range *s {
		x = append(x, player.PlayerID)
	}

	m.PlayerIDs = x
	m.UpdatedAt = f.Clock.Now()

	return m
}
