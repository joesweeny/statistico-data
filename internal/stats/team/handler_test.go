package team_stats

import (
	"testing"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHandleTeamStats(t *testing.T) {
	t.Run("returns a slice of proto TeamStats structs", func(t *testing.T) {
		x := []*model.TeamStats{
			modelTeamStats(),
			modelTeamStats(),
			modelTeamStats(),
		}

		stats := HandleTeamStats(x)

		assert.Equal(t, 3, len(stats))
	})
}

func modelTeamStats() *model.TeamStats {
	total := 5
	blocked := 2
	corners := 15
	redCards := 1
	dangerous := 30
	goalKicks := 12

	return &model.TeamStats{
		TeamID: 850,
		TeamShots: model.TeamShots{
			Total: &total,
			Blocked: &blocked,
		},
		Corners: &corners,
		RedCards: &redCards,
		TeamAttacks: model.TeamAttacks{
			Dangerous: &dangerous,
		},
		GoalKicks: &goalKicks,
	}
}