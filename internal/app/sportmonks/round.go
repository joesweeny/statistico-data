package sportmonks

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

type RoundRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (r RoundRequester) RoundsBySeasonIDs(seasonIDs []uint64) <-chan *app.Round {
	ch := make(chan *app.Round, 500)

	go r.parseRounds(seasonIDs, ch)

	return ch
}

func (r RoundRequester) parseRounds(seasonIDs []uint64, ch chan<- *app.Round) {
	defer close(ch)

	for _, id := range seasonIDs {
		r.sendRoundRequest(id, ch)
	}
}

func (r RoundRequester) sendRoundRequest(seasonID uint64, ch chan<- *app.Round) {
	res, _, err := r.client.RoundsBySeasonID(context.Background(), int(seasonID), []string{})

	if err != nil {
		r.logger.Errorf("Error when calling client '%s' when making round request", err.Error())
		return
	}

	for _, round := range res {
		x, err := transformRound(&round)

		if err != nil {
			r.logger.Warningf(err.Error())
			continue
		}

		ch <- x
	}
}

func transformRound(r *spClient.Round) (*app.Round, error) {
	start, err := time.Parse(dateFormat, r.Start)

	if err != nil {
		return nil, fmt.Errorf("error parsing round from client. ID '%d', error %s", r.ID, err)
	}

	end, err := time.Parse(dateFormat, r.End)

	if err != nil {
		return nil, fmt.Errorf("error parsing round from client. ID '%d', error %s", r.ID, err)
	}

	return &app.Round{
		ID:        uint64(r.ID),
		Name:      strconv.Itoa(r.Name),
		SeasonID:  uint64(r.SeasonID),
		StartDate: start,
		EndDate:   end,
	}, nil
}

func NewRoundRequester(client *spClient.HTTPClient, log *logrus.Logger) *RoundRequester {
	return &RoundRequester{client: client, logger: log}
}
