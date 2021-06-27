package sportmonks

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-football-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type PlayerRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (p PlayerRequester) PlayerByID(id uint64) (*app.Player, error) {
	res, _, err := p.client.PlayerByID(context.Background(), int(id), []string{})

	if err != nil {
		p.logger.Errorf("Error calling client '%s' when making player request, player id %ds", err.Error(), id)
		return nil, fmt.Errorf("unable to fetch player with id %d", id)
	}

	return transformPlayer(res), nil
}

func transformPlayer(p *spClient.Player) *app.Player {
	return &app.Player{
		ID:          uint64(p.ID),
		CountryId:   uint64(p.CountryID),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		BirthPlace:  &p.BirthPlace,
		DateOfBirth: &p.BirthDate,
		PositionID:  p.PositionID,
		Image:       p.ImagePath,
	}
}

func NewPlayerRequester(client *spClient.HTTPClient, log *logrus.Logger) *PlayerRequester {
	return &PlayerRequester{client: client, logger: log}
}
