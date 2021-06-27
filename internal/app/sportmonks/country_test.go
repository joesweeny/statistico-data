package sportmonks_test

import (
	"bytes"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-football-data/internal/app/mock"
	"github.com/statistico/statistico-football-data/internal/app/sportmonks"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCountries(t *testing.T) {
	t.Run("countries returns channel", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(countriesResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewCountryRequester(&client, logger)

		ch := requester.Countries()

		eng := <-ch
		ger := <-ch

		a := assert.New(t)

		a.Equal(uint64(180), eng.ID)
		a.Equal("England", eng.Name)
		a.Equal("Europe", eng.Continent)
		a.Equal("ENG", eng.ISO)
		a.Equal(uint64(5), ger.ID)
		a.Equal("Germany", ger.Name)
		a.Equal("Europe", ger.Continent)
		a.Equal("DEU", ger.ISO)
	})
}

var countriesResponse = `{
	"data": [
		{
			"id": 180,
			"name": "England",
			"extra": {
				"continent": "Europe",
				"sub_region": "Western Europe",
				"world_region": "Europe",
				"fifa": "GER",
				"iso": "ENG",
				"iso2": "EN"
			}
		},
		{
			"id": 5,
			"name": "Germany",
			"extra": {
				"continent": "Europe",
				"sub_region": "Western Europe",
				"world_region": "Europe",
				"fifa": "GER",
				"iso": "DEU",
				"iso2": "DE"
			}
		}
	],
	"meta": {
		"pagination": {
			"total": 2,
			"count": 2,
			"per_page": 100,
			"current_page": 1,
			"total_pages": 1,
			"links": {}
		}
	}
}`
