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

func TestSeasonRequester_Seasons(t *testing.T) {
	t.Run("returns a channel containing season struct", func(t *testing.T) {
		t.Helper()

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(seasonsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewSeasonRequester(&client, logger)

		ch := requester.Seasons()

		season := <-ch

		a := assert.New(t)

		a.Equal(uint64(16029), season.ID)
		a.Equal("2019/2020", season.Name)
		a.Equal(uint64(2), season.CompetitionID)
		a.Equal(true, season.IsCurrent)
	})
}

var seasonsResponse = `{
	"data": [
		{
			"id": 16029,
			"name": "2019/2020",
			"league_id": 2,
			"is_current_season": true,
			"current_round_id": 183973,
			"current_stage_id": 77443828
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
