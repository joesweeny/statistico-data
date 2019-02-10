package stats

import (
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
)

type TeamFactory struct {
	Clock clockwork.Clock
}

func (f TeamFactory) createTeamStats(s *sportmonks.TeamStats) *model.TeamStats {
	shots := model.TeamShots{
		Total:      s.Shots.Total,
		OnGoal:     s.Shots.Ongoal,
		OffGoal:    s.Shots.Offgoal,
		Blocked:    s.Shots.Blocked,
		InsideBox:  s.Shots.Insidebox,
		OutsideBox: s.Shots.Outsidebox,
	}

	passes := model.TeamPasses{
		Total:      s.Passes.Total,
		Accuracy:   s.Passes.Accurate,
		Percentage: s.Passes.Percentage,
	}

	attacks := model.TeamAttacks{
		Total:     s.Attacks.Attacks,
		Dangerous: s.Attacks.DangerousAttacks,
	}
	return &model.TeamStats{
		FixtureID:     s.FixtureID,
		TeamID:        s.TeamID,
		TeamShots:     shots,
		TeamPasses:    passes,
		TeamAttacks:   attacks,
		Fouls:         s.Fouls,
		Corners:       s.Corners,
		Offsides:      s.Offsides,
		Possession:    s.Possessiontime,
		YellowCards:   s.Yellowcards,
		RedCards:      s.Redcards,
		Saves:         s.Saves,
		Substitutions: s.Substitutions,
		GoalKicks:     s.GoalKick,
		GoalAttempts:  s.GoalAttempts,
		FreeKicks:     s.FreeKick,
		ThrowIns:      s.ThrowIn,
		CreatedAt:     f.Clock.Now(),
		UpdatedAt:     f.Clock.Now(),
	}
}