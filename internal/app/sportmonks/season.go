package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-football-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type SeasonRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (s SeasonRequester) Seasons() <-chan *app.Season {
	_, meta, err := s.client.Seasons(context.Background(), 1, []string{})

	if err != nil {
		s.logger.Errorf("Error when calling client '%s' when making season request", err.Error())
		return nil
	}

	ch := make(chan *app.Season, meta.Pagination.Total)

	go s.parseSeasons(meta.Pagination.TotalPages, ch)

	return ch
}

func (s SeasonRequester) parseSeasons(pages int, ch chan<- *app.Season) {
	defer close(ch)

	for i := 1; i <= pages; i++ {
		s.sendSeasonRequest(i, ch)
	}
}

func (s SeasonRequester) sendSeasonRequest(page int, ch chan<- *app.Season) {
	res, _, err := s.client.Seasons(context.Background(), page, []string{})

	if err != nil {
		s.logger.Errorf("Error when calling client '%s' when making season request", err.Error())
		return
	}

	for _, season := range res {
		ch <- transformSeason(&season)
	}
}

func transformSeason(s *spClient.Season) *app.Season {
	return &app.Season{
		ID:            uint64(s.ID),
		Name:          s.Name,
		CompetitionID: uint64(s.LeagueID),
		IsCurrent:     s.IsCurrentSeason,
	}
}

func NewSeasonRequester(client *spClient.HTTPClient, log *logrus.Logger) *SeasonRequester {
	return &SeasonRequester{client: client, logger: log}
}
