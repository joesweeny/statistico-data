package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type CompetitionRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (c CompetitionRequester) Competitions() <-chan *app.Competition {
	_, meta, err := c.client.Leagues(context.Background(), 1, []string{})

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s' when making competition request", err.Error())
		return nil
	}

	ch := make(chan *app.Competition, meta.Pagination.Total)

	go c.parseCompetitions(meta.Pagination.TotalPages, ch)

	return ch
}

func (c CompetitionRequester) parseCompetitions(pages int, ch chan<- *app.Competition) {
	defer close(ch)

	for i := 1; i <= pages; i++ {
		c.sendCompetitionRequest(i, ch)
	}
}

func (c CompetitionRequester) sendCompetitionRequest(page int, ch chan<- *app.Competition) {
	res, _, err := c.client.Leagues(context.Background(), page, []string{})

	if err != nil {
		c.logger.Fatalf("Error when calling client '%s' when making competition request", err.Error())
		return
	}

	for _, competition := range res {
		ch <- transformCompetition(&competition)
	}
}

func transformCompetition(s *spClient.League) *app.Competition {
	return &app.Competition{
		ID:        uint64(s.ID),
		Name:      s.Name,
		CountryID: uint64(s.CountryID),
		IsCup:     s.IsCup,
	}
}

func NewCompetitionRequester(client *spClient.HTTPClient, log *logrus.Logger) *CompetitionRequester {
	return &CompetitionRequester{client: client, logger: log}
}
