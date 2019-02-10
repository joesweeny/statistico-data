package event

import (
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"strconv"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createGoalEvent(s *sportmonks.GoalEvent) *model.GoalEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return &model.GoalEvent{
		ID:             s.ID,
		FixtureID:      s.FixtureID,
		TeamID:         teamId,
		PlayerID:       s.PlayerID,
		PlayerAssistID: s.PlayerAssistID,
		Minute:         s.Minute,
		Score:          s.Result,
		CreatedAt:      f.Clock.Now(),
	}
}

func (f Factory) createSubstitutionEvent(s *sportmonks.SubstitutionEvent) *model.SubstitutionEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return &model.SubstitutionEvent{
		ID:          s.ID,
		FixtureID:   s.FixtureID,
		TeamID:      teamId,
		PlayerInID:  s.PlayerInID,
		PlayerOutID: s.PlayerOutID,
		Minute:      s.Minute,
		Injured:     s.Injured,
		CreatedAt:   f.Clock.Now(),
	}
}
