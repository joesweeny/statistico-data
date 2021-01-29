package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"sync"
)

type ResultRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (r ResultRequester) ResultsBySeasonIDs(seasonIDs []uint64) <-chan app.Result {
	ch := make(chan app.Result, 100)

	go r.parseResults(seasonIDs, ch)

	return ch
}

func (r ResultRequester) parseResults(seasonIDs []uint64, ch chan<- app.Result) {
	defer close(ch)

	var wg sync.WaitGroup

	for _, id := range seasonIDs {
		wg.Add(1)
		go r.sendSeasonRequests(id, ch, &wg)
	}

	wg.Wait()
}

func (r ResultRequester) sendSeasonRequests(seasonID uint64, ch chan<- app.Result, w *sync.WaitGroup) {
	season, _, err := r.client.SeasonByID(context.Background(), int(seasonID), []string{"results"})

	if err != nil {
		r.logger.Errorf(
			"Error when calling client '%s' when making fixtures request. Season ID %d",
			err.Error(),
			seasonID,
		)

		w.Done()
		return
	}

	for _, result := range season.Results() {
		ch <- transformResult(result)
	}

	w.Done()
}

func transformResult(s spClient.Fixture) app.Result {
	return app.Result{
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
