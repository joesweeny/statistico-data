package result

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createResult(s *sportmonks.Fixture) *model.Result {
	return &model.Result{
		FixtureID:          s.ID,
		PitchCondition:     s.Pitch,
		HomeFormation:      &s.Formations.LocalteamFormation,
		AwayFormation:      &s.Formations.VisitorteamFormation,
		HomeScore:          &s.Scores.LocalteamScore,
		AwayScore:          &s.Scores.VisitorteamScore,
		HomePenScore:       s.Scores.LocalteamPenScore,
		AwayPenScore:       s.Scores.VisitorteamPenScore,
		HalfTimeScore:      s.Scores.HtScore,
		FullTimeScore:      s.Scores.FtScore,
		ExtraTimeScore:     s.Scores.EtScore,
		HomeLeaguePosition: &s.Standings.LocalteamPosition,
		AwayLeaguePosition: &s.Standings.VisitorteamPosition,
		Minutes:            &s.Time.Minute,
		Seconds:            s.Time.Second,
		AddedTime:          s.Time.AddedTime,
		ExtraTime:          s.Time.ExtraMinute,
		InjuryTime:         s.Time.InjuryTime,
		CreatedAt:          f.Clock.Now(),
		UpdatedAt:          f.Clock.Now(),
	}
}
