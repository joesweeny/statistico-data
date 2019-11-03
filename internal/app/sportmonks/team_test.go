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

func TestTeamsBySeasonIDs(t *testing.T) {
	t.Run("returns channel containing team struct", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(teamsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewTeamRequester(&client, logger)

		ch := requester.TeamsBySeasonIDs([]uint64{uint64(435), uint64(33), uint64(2)})

		x := <-ch
		y := <-ch
		z := <-ch

		a := assert.New(t)

		a.Equal(uint64(1), x.ID)
		a.Equal("West Ham United", x.Name)
		a.Equal("WHU", *x.ShortCode)
		a.Equal(uint64(462), x.CountryID)
		a.Equal(uint64(214), x.VenueID)
		a.Equal(false, x.NationalTeam)
		a.Equal(1895, *x.Founded)
		a.Equal("https://cdn.sportmonks.com/images/soccer/teams/1/1.png", *x.Logo)

		a.Equal(uint64(1), y.ID)
		a.Equal("West Ham United", y.Name)
		a.Equal("WHU", *y.ShortCode)
		a.Equal(uint64(462), y.CountryID)
		a.Equal(uint64(214), y.VenueID)
		a.Equal(false, y.NationalTeam)
		a.Equal(1895, *y.Founded)
		a.Equal("https://cdn.sportmonks.com/images/soccer/teams/1/1.png", *y.Logo)

		a.Equal(uint64(1), z.ID)
		a.Equal("West Ham United", z.Name)
		a.Equal("WHU", *z.ShortCode)
		a.Equal(uint64(462), z.CountryID)
		a.Equal(uint64(214), z.VenueID)
		a.Equal(false, z.NationalTeam)
		a.Equal(1895, *z.Founded)
		a.Equal("https://cdn.sportmonks.com/images/soccer/teams/1/1.png", *z.Logo)
	})
}

var teamsResponse = `{
	"data": [
		{
			"id": 1,
			"legacy_id": 377,
			"name": "West Ham United",
			"short_code": "WHU",
			"twitter": "@WestHamUtd",
			"country_id": 462,
			"national_team": false,
			"founded": 1895,
			"logo_path": "https:\/\/cdn.sportmonks.com\/images\/soccer\/teams\/1\/1.png",
			"venue_id": 214,
			"current_season_id": 16036
		}
	],
	"meta": {
		"pagination": {
			"total": 1,
			"count": 1,
			"per_page": 100,
			"current_page": 1,
			"total_pages": 1,
			"links": {}
		}
	}
}`
