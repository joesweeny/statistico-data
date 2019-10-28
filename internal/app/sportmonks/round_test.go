package sportmonks_test

import (
	"bytes"
	"github.com/sirupsen/logrus"

	//"github.com/sirupsen/logrus"

	//"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRoundsBySeasonIDs(t *testing.T) {
	t.Run("returns channel containing round struct", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(roundsResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewRoundRequester(&client, logger)

		ch := requester.RoundsBySeasonIDs([]int64{int64(100), int64(234)})

		x := <- ch
		y := <- ch

		a := assert.New(t)

		a.Equal(int64(100), x.ID)
		a.Equal("5", x.Name)
		a.Equal(int64(9), x.SeasonID)
		a.Equal("2011-09-17", x.StartDate.Format("2006-01-02"))
		a.Equal("2011-09-18", x.EndDate.Format("2006-01-02"))

		a.Equal(int64(234), y.ID)
		a.Equal("34", y.Name)
		a.Equal(int64(9), y.SeasonID)
		a.Equal("2018-09-17", y.StartDate.Format("2006-01-02"))
		a.Equal("2019-09-18", y.EndDate.Format("2006-01-02"))
	})

	t.Run("error is logged if unable to parse date", func(t *testing.T) {
		res := `{
			"data": [
				{
					"id": 100,
					"name": 5,
					"league_id": 8,
					"season_id": 9,
					"stage_id": 7,
					"start": "2011-09-17",
					"end": "2011-09-18"
				},
				{
					"id": 234,
					"name": 34,
					"league_id": 8,
					"season_id": 9,
					"stage_id": 7,
					"start": "Today",
					"end": "2019-09-18"
				}
			]
		}`

		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(res)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, hook := test.NewNullLogger()

		requester := sportmonks.NewRoundRequester(&client, logger)

		ch := requester.RoundsBySeasonIDs([]int64{int64(100), int64(234)})

		x := <- ch

		a := assert.New(t)

		a.Equal(int64(100), x.ID)
		a.Equal("5", x.Name)
		a.Equal(int64(9), x.SeasonID)
		a.Equal("2011-09-17", x.StartDate.Format("2006-01-02"))
		a.Equal("2011-09-18", x.EndDate.Format("2006-01-02"))

		t.Errorf("Entry: %+v", hook.LastEntry().Level)

		//a.Equal(t, int(1), len(hook.Entries))
		a.Equal(t, logrus.WarnLevel.String(), hook.LastEntry().Level)
		a.Equal(
			`error parsing round from client. ID '234', error parsing time "Today" as "2006-01-02": cannot parse "Today" as "2006"`,
			hook.LastEntry().Message,
		)
	})
}

var roundsResponse = `{
	"data": [
		{
			"id": 100,
			"name": 5,
			"league_id": 8,
			"season_id": 9,
			"stage_id": 7,
			"start": "2011-09-17",
			"end": "2011-09-18"
		},
		{
			"id": 234,
			"name": 34,
			"league_id": 8,
			"season_id": 9,
			"stage_id": 7,
			"start": "2018-09-17",
			"end": "2019-09-18"
		}
	]
}`
