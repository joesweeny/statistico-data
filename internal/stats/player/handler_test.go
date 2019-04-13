package player_stats

import (
	"testing"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHandlePlayerStats(t *testing.T) {
	t.Run("returns a slice of proto PlayerStats structs", func(t *testing.T) {
		var (
			goals = 1
			assists = 1
			onGoal = 2
		)

		x := []model.PlayerStats{
			*modelPlayerStats(&goals, &assists, &onGoal),
			*modelPlayerStats(&goals, &assists, &onGoal),
			*modelPlayerStats(&goals, &assists, &onGoal),
		}

		stats := HandlePlayerStats(&x)

		assert.Equal(t, 3, len(stats))
	})
}

func modelPlayerStats(goals *int, assists *int, onGoal *int) *model.PlayerStats {
	shots := 5
	conceded := 0
	return &model.PlayerStats{
		PlayerID:        77,
		PlayerShots:     model.PlayerShots{
			Total: 	&shots,
			OnGoal: onGoal,
		},
		PlayerGoals:     model.PlayerGoals{
			Scored: goals,
			Conceded: &conceded,
		},
		Assists: assists,
	}
}