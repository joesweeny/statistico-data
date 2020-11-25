package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type VenueRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (v VenueRequester) VenuesBySeasonIDs(seasonIDs []uint64) <-chan *app.Venue {
	ch := make(chan *app.Venue, 500)

	go v.parseVenues(seasonIDs, ch)

	return ch
}

func (v VenueRequester) parseVenues(seasonIDs []uint64, ch chan<- *app.Venue) {
	defer close(ch)

	for _, id := range seasonIDs {
		v.callClient(id, ch)
	}
}

func (v VenueRequester) callClient(seasonId uint64, ch chan<- *app.Venue) {
	res, _, err := v.client.VenuesBySeasonID(context.Background(), int(seasonId))

	if err != nil {
		v.logger.Errorf("Error when calling client '%s' when making venue request", err.Error())
		return
	}

	for _, venue := range res {
		ch <- transformVenue(&venue)
	}
}

func transformVenue(v *spClient.Venue) *app.Venue {
	return &app.Venue{
		ID:       uint64(v.ID),
		Name:     v.Name,
		Surface:  &v.Surface,
		Address:  v.Address,
		City:     &v.City,
		Capacity: &v.Capacity,
	}
}

func NewVenueRequester(client *spClient.HTTPClient, log *logrus.Logger) *VenueRequester {
	return &VenueRequester{client: client, logger: log}
}
