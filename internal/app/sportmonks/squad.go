package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"sync"
)

type SquadRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (s SquadRequester) SquadsBySeasonIDs(seasonIDs []uint64) <-chan *app.Squad {
	ch := make(chan *app.Squad, 500)

	go s.parseSquads(seasonIDs, ch)

	return ch
}

func (s SquadRequester) parseSquads(seasonIDs []uint64, ch chan<- *app.Squad) {
	defer close(ch)

	var wg sync.WaitGroup

	for _, id := range seasonIDs {
		wg.Add(1)
		go s.sendSquadRequests(id, ch, &wg)
	}

	wg.Wait()
}

func (s SquadRequester) sendSquadRequests(seasonID uint64, ch chan<- *app.Squad, w *sync.WaitGroup) {
	_, meta, err := s.client.TeamsBySeasonID(context.Background(), int(seasonID), 1, []string{"squad"})

	if err != nil {
		s.logger.Fatalf("Error when calling client '%s' when making squad request", err.Error())
		return
	}

	for i := 1; i <= meta.Pagination.TotalPages; i++ {
		res, _, err := s.client.TeamsBySeasonID(context.Background(), int(seasonID), i, []string{"squad"})

		if err != nil {
			s.logger.Fatalf("Error when calling client '%s' when making squad request", err.Error())
			return
		}

		for _, team := range res {
			ch <- transformSquad(&team, seasonID)
		}
	}

	w.Done()
}

func transformSquad(t *spClient.Team, seasonID uint64) *app.Squad {
	squad := app.Squad{
		SeasonID: seasonID,
		TeamID:   uint64(t.ID),
	}

	for _, player := range t.Squad() {
		squad.PlayerIDs = append(squad.PlayerIDs, uint64(player.PlayerID))
	}

	return &squad
}

func NewSquadRequester(client *spClient.HTTPClient, log *logrus.Logger) *SquadRequester {
	return &SquadRequester{client: client, logger: log}
}
