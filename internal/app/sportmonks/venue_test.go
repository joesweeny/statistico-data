package sportmonks_test

import (
	"bytes"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestVenuesBySeasonIDs(t *testing.T) {
	t.Run("returns channel containing venue struct", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(venuesResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewVenueRequester(&client, logger)

		ch := requester.VenuesBySeasonIDs([]uint64{uint64(500)})

		ven := <-ch

		a := assert.New(t)

		a.Equal(uint64(200), ven.ID)
		a.Equal("Turf Moor", ven.Name)
		a.Equal("grass", *ven.Surface)
		a.Equal("Harry Potts Way", *ven.Address)
		a.Equal("Burnley", *ven.City)
		a.Equal(22546, *ven.Capacity)
	})
}

var venuesResponse = `{
	"data": [
		{
			"id": 200,
			"name": "Turf Moor",
			"surface": "grass",
			"address": "Harry Potts Way",
			"city": "Burnley",
			"capacity": 22546,
			"image_path": "https:\/\/cdn.sportmonks.com\/images\/soccer\/venues\/8\/200.png",
			"coordinates": null
		}
	]
}`
