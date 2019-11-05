package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type ResultRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (r ResultRequester) ResultByFixtureID(id uint64) (*app.Result, error) {
	var filters map[string][]int
	var includes []string

	res, _ , err := r.client.FixtureByID(context.Background(), int(id), includes, filters)

	if err != nil {
		return nil, err
	}

	return transformResult(res), nil
}

func transformResult(s *spClient.Fixture) *app.Result {
	return &app.Result{
		FixtureID:          uint64(s.ID),
		PitchCondition:     s.Pitch,
		HomeFormation:      &s.Formations.LocalTeamFormation,
		AwayFormation:      &s.Formations.VisitorTeamFormation,
		HomeScore:          &s.Scores.LocalTeamScore,
		AwayScore:          &s.Scores.VisitorTeamScore,
		HomePenScore:       s.Scores.LocalTeamPenScore,
		AwayPenScore:       s.Scores.VisitorTeamPenScore,
		HalfTimeScore:      s.Scores.HTScore,
		FullTimeScore:      s.Scores.FTScore,
		ExtraTimeScore:     s.Scores.ETScore,
		HomeLeaguePosition: &s.Standings.LocalTeamPosition,
		AwayLeaguePosition: &s.Standings.VisitorTeamPosition,
		Minutes:            &s.Time.Minute,
		AddedTime:          s.Time.AddedTime,
		ExtraTime:          s.Time.ExtraMinute,
		InjuryTime:         s.Time.InjuryTime,
	}
}

func NewResultRequester(client *spClient.HTTPClient, log *logrus.Logger) *ResultRequester {
	return &ResultRequester{client: client, logger: log}
}
