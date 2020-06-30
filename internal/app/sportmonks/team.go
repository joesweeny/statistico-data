package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"sync"
)

type TeamRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (t TeamRequester) TeamsBySeasonIDs(seasonIDs []uint64) <-chan *app.Team {
	ch := make(chan *app.Team, 500)

	go t.parseTeams(seasonIDs, ch)

	return ch
}

func (t TeamRequester) parseTeams(seasonIDs []uint64, ch chan<- *app.Team) {
	defer close(ch)

	var waitGroup sync.WaitGroup

	for _, id := range seasonIDs {
		waitGroup.Add(1)
		go t.sendTeamRequests(id, ch, &waitGroup)
	}

	waitGroup.Wait()
}

func (t TeamRequester) sendTeamRequests(seasonID uint64, ch chan<- *app.Team, w *sync.WaitGroup) {
	_, meta, err := t.client.TeamsBySeasonID(context.Background(), int(seasonID), 1, []string{})

	if err != nil {
		t.logger.Fatalf("Error when calling client '%s' when making team request", err.Error())
		return
	}

	for i := 1; i <= meta.Pagination.TotalPages; i++ {
		res, _, err := t.client.TeamsBySeasonID(context.Background(), int(seasonID), i, []string{})

		if err != nil {
			t.logger.Fatalf("Error when calling client '%s' when making team request", err.Error())
			return
		}

		for _, team := range res {
			ch <- transformTeam(team)
		}
	}

	w.Done()
}

func transformTeam(t spClient.Team) *app.Team {
	return &app.Team{
		ID:           uint64(t.ID),
		Name:         t.Name,
		ShortCode:    &t.ShortCode,
		CountryID:    uint64(t.CountryID),
		VenueID:      uint64(t.VenueID),
		NationalTeam: t.NationalTeam,
		Founded:      &t.Founded,
		Logo:         t.LogoPath,
	}
}

func NewTeamRequester(client *spClient.HTTPClient, log *logrus.Logger) *TeamRequester {
	return &TeamRequester{client: client, logger: log}
}
