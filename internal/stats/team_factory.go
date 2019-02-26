package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statistico-data/internal/model"
	"github.com/jonboulle/clockwork"
	"strconv"
)

type TeamFactory struct {
	Clock clockwork.Clock
}

func (f TeamFactory) createTeamStats(s *sportmonks.TeamStats) *model.TeamStats {
	return &model.TeamStats{
		FixtureID:     s.FixtureID,
		TeamID:        s.TeamID,
		TeamShots:     *handleTeamShots(&s.Shots),
		TeamPasses:    *handleTeamPasses(&s.Passes),
		TeamAttacks:   *handleTeamAttacks(&s.Attacks),
		Fouls:         parseInt(s.Fouls),
		Corners:       parseInt(s.Corners),
		Offsides:      parseInt(s.Offsides),
		Possession:    parseInt(s.Possessiontime),
		YellowCards:   parseInt(s.Yellowcards),
		RedCards:      parseInt(s.Redcards),
		Saves:         parseInt(s.Saves),
		Substitutions: parseInt(s.Substitutions),
		GoalKicks:     parseInt(s.GoalKick),
		GoalAttempts:  parseInt(s.GoalAttempts),
		FreeKicks:     parseInt(s.FreeKick),
		ThrowIns:      parseInt(s.ThrowIn),
		CreatedAt:     f.Clock.Now(),
		UpdatedAt:     f.Clock.Now(),
	}
}

func (f TeamFactory) updateTeamStats(s *sportmonks.TeamStats, m *model.TeamStats) *model.TeamStats {
	m.TeamShots = *handleTeamShots(&s.Shots)
	m.TeamPasses = *handleTeamPasses(&s.Passes)
	m.TeamAttacks = *handleTeamAttacks(&s.Attacks)
	m.Fouls = parseInt(s.Fouls)
	m.Corners = parseInt(s.Corners)
	m.Offsides = parseInt(s.Offsides)
	m.Possession = parseInt(s.Possessiontime)
	m.YellowCards = parseInt(s.Yellowcards)
	m.RedCards = parseInt(s.Redcards)
	m.Saves = parseInt(s.Saves)
	m.Substitutions = parseInt(s.Substitutions)
	m.GoalKicks = parseInt(s.GoalKick)
	m.GoalAttempts = parseInt(s.GoalAttempts)
	m.FreeKicks = parseInt(s.FreeKick)
	m.ThrowIns = parseInt(s.ThrowIn)
	m.UpdatedAt = f.Clock.Now()

	return m
}

func handleTeamShots(s *sportmonks.TeamShots) *model.TeamShots {
	return &model.TeamShots{
		Total:      parseInt(s.Total),
		OnGoal:     parseInt(s.Ongoal),
		OffGoal:    parseInt(s.Offgoal),
		Blocked:    parseInt(s.Blocked),
		InsideBox:  parseInt(s.Insidebox),
		OutsideBox: parseInt(s.Outsidebox),
	}
}

func handleTeamPasses(s *sportmonks.TeamPasses) *model.TeamPasses {
	return &model.TeamPasses{
		Total:      parseInt(s.Total),
		Accuracy:   parseInt(s.Accurate),
		Percentage: parseInt(s.Percentage),
	}
}

func handleTeamAttacks(s *sportmonks.TeamAttacks) *model.TeamAttacks {
	return &model.TeamAttacks{
		Total:     parseInt(s.Attacks),
		Dangerous: parseInt(s.DangerousAttacks),
	}
}

// Some stats are being sent as either int or string, this function here is a helper
// to ensure the property value is consistent as an int
func parseInt(i interface{}) *int {
	if i == nil {
		return nil
	}

	if _, ok := i.(int); ok {
		val := i.(int)
		return &val
	}

	if x, ok := i.(float64); ok {
		val := int(x)
		return &val
	}

	if _, ok := i.(string); ok {
		val, err := strconv.Atoi(i.(string))

		if err != nil {
			panic(err)
		}

		return &val
	}

	return nil
}
