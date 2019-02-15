package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
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

func (f TeamFactory) updateTeamStats(s *sportmonks.TeamStats, m *model.TeamStats) *model.TeamStats {
	m.TeamShots = *handleTeamShots(&s.Shots)
	m.TeamPasses = *handleTeamPasses(&s.Passes)
	m.TeamAttacks = *handleTeamAttacks(&s.Attacks)
	m.Fouls = s.Fouls
	m.Corners = s.Corners
	m.Offsides = s.Offsides
	m.Possession = s.Possessiontime
	m.YellowCards = s.Yellowcards
	m.RedCards = s.Redcards
	m.Saves = s.Saves
	m.Substitutions = s.Substitutions
	m.GoalKicks = s.GoalKick
	m.GoalAttempts = s.GoalAttempts
	m.FreeKicks = s.FreeKick
	m.ThrowIns = s.ThrowIn
	m.UpdatedAt = f.Clock.Now()

	return m
}

func handleTeamShots(s *sportmonks.TeamShots) *model.TeamShots {
	return &model.TeamShots{
		Total:      s.Total,
		OnGoal:     s.Ongoal,
		OffGoal:    s.Offgoal,
		Blocked:    s.Blocked,
		InsideBox:  s.Insidebox,
		OutsideBox: s.Outsidebox,
	}
}

func handleTeamPasses(s *sportmonks.TeamPasses) *model.TeamPasses {
	return &model.TeamPasses{
		Total:      s.Total,
		Accuracy:   s.Accurate,
		Percentage: s.Percentage,
	}
}

func handleTeamAttacks(s *sportmonks.TeamAttacks) *model.TeamAttacks {
	t := parseInt(s.Attacks)
	d := parseInt(s.DangerousAttacks)
	return &model.TeamAttacks{
		Total:     &t,
		Dangerous: &d,
	}
}

// Some stats are being sent as either int or string, this function here is a helper
// to ensure the property value is consistent as an int
func parseInt(i interface{}) int {
	_, ok := i.(int)

	if ok {
		return i.(int)
	}

	x, err := strconv.Atoi(i.(string))

	if err != nil {
		panic(err)
	}

	return x
}