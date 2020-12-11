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

func (r ResultRequester) ResultsByFixtureIDs(ids []uint64) <-chan *app.Result {
	ch := make(chan *app.Result, 100)

	go r.parseResults(ids, ch)

	return ch
}

func (r ResultRequester) parseResults(ids []uint64, ch chan<- *app.Result) {
	defer close(ch)

	var filters map[string][]int
	var includes []string

	for _, id := range ids {
		res, _, err := r.client.FixtureByID(context.Background(), int(id), includes, filters)

		if err != nil {
			r.logger.Errorf(
				"Error when calling client '%s' when making requests request. Fixture ID %d",
				err.Error(),
				id,
			)
			return
		}

		ch <- transformResult(res)
	}
}

func transformResult(s *spClient.Fixture) *app.Result {
	return &app.Result{
		FixtureID:          uint64(s.ID),
		PitchCondition:     s.Pitch,
		HomeFormation:      s.Formations.LocalTeamFormation,
		AwayFormation:      s.Formations.VisitorTeamFormation,
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
