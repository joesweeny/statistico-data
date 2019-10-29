package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type PlayerRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (p PlayerRequester) PlayerByID(id int64) *app.Player {
	res, _, err := p.client.PlayerByID(context.Background(), int(id), []string{})

	if err != nil {
		p.logger.Fatalf("Error calling client '%s' when making player request", err.Error())
		return nil
	}

	return transformPlayer(res)
}

func transformPlayer(p *spClient.Player) *app.Player {
	return &app.Player{
		ID:          int64(p.ID),
		CountryId:   int64(p.CountryID),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		BirthPlace:  &p.BirthPlace,
		DateOfBirth: &p.BirthDate,
		PositionID:  p.PositionID,
		Image:       &p.ImagePath,
	}
}

func NewPlayerRequester(client *spClient.HTTPClient, log *logrus.Logger) *PlayerRequester {
	return &PlayerRequester{client: client, logger: log}
}
