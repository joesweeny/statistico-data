package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/helpers"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"sync"
	"time"
)

type FixtureRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (f FixtureRequester) FixturesBySeasonIDs(ids []uint64) <-chan *app.Fixture {
	ch := make(chan *app.Fixture, 100)

	go f.parseFixtures(ids, ch)

	return ch
}

func (f FixtureRequester) parseFixtures(seasonIDs []uint64, ch chan<- *app.Fixture) {
	defer close(ch)

	var wg sync.WaitGroup

	for _, id := range seasonIDs {
		wg.Add(1)
		go f.sendFixtureRequests(id, ch, &wg)
	}

	wg.Wait()
}

func (f FixtureRequester) sendFixtureRequests(seasonID uint64, ch chan<- *app.Fixture, w *sync.WaitGroup) {
	res, _, err := f.client.SeasonByID(context.Background(), int(seasonID), []string{"fixtures"})

	if err != nil {
		f.logger.Warnf(
			"Error when calling client '%s' when making fixtures request. Season ID %d",
			err.Error(),
			seasonID,
		)

		w.Done()
		return
	}

	for _, fixture := range res.Fixtures() {
		ch <- transformFixture(&fixture)
	}

	w.Done()
}

func transformFixture(s *spClient.Fixture) *app.Fixture {
	return &app.Fixture{
		ID:         uint64(s.ID),
		SeasonID:   uint64(s.SeasonID),
		RoundID:    helpers.NullableUint64(s.RoundID),
		VenueID:    helpers.NullableUint64(s.VenueID),
		HomeTeamID: uint64(s.LocalTeamID),
		AwayTeamID: uint64(s.VisitorTeamID),
		RefereeID:  helpers.NullableUint64(s.RefereeID),
		Date:       time.Unix(int64(s.Time.StartingAt.Timestamp), 0),
	}
}

func NewFixtureRequester(client *spClient.HTTPClient, log *logrus.Logger) *FixtureRequester {
	return &FixtureRequester{client: client, logger: log}
}
