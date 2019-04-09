package result

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
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

func (f Factory) updateResult(s *sportmonks.Fixture, m *model.Result) *model.Result {
	m.PitchCondition = s.Pitch
	m.HomeFormation = &s.Formations.LocalteamFormation
	m.AwayFormation = &s.Formations.VisitorteamFormation
	m.HomeScore = &s.Scores.LocalteamScore
	m.AwayScore = &s.Scores.VisitorteamScore
	m.HomePenScore = s.Scores.LocalteamPenScore
	m.AwayPenScore = s.Scores.VisitorteamPenScore
	m.HalfTimeScore = s.Scores.HtScore
	m.FullTimeScore = s.Scores.FtScore
	m.ExtraTimeScore = s.Scores.EtScore
	m.HomeLeaguePosition = &s.Standings.LocalteamPosition
	m.AwayLeaguePosition = &s.Standings.VisitorteamPosition
	m.Minutes = &s.Time.Minute
	m.Seconds = s.Time.Second
	m.AddedTime = s.Time.AddedTime
	m.ExtraTime = s.Time.ExtraMinute
	m.InjuryTime = s.Time.InjuryTime
	m.UpdatedAt = f.Clock.Now()

	return m
}
