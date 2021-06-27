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

func TestPlayerRequester_PlayerByID(t *testing.T) {
	t.Run("returns player struct", func(t *testing.T) {
		server := mock.HttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(playerResponse)),
			}, nil
		})

		client := spClient.HTTPClient{
			HTTPClient: server,
			BaseURL:    "http://example.com",
			Key:        "my-key",
		}

		logger, _ := test.NewNullLogger()

		requester := sportmonks.NewPlayerRequester(&client, logger)

		player, err := requester.PlayerByID(uint64(219591))

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		a := assert.New(t)

		a.Equal(uint64(219591), player.ID)
		a.Equal(uint64(1190), player.CountryId)
		a.Equal("Fabián Cornelio", player.FirstName)
		a.Equal("Balbuena González", player.LastName)
		a.Equal("Ciudad del Este", *player.BirthPlace)
		a.Equal("23/08/1991", *player.DateOfBirth)
		a.Equal(2, player.PositionID)
		a.Equal("https://cdn.sportmonks.com/images/soccer/player/1/1.png", player.Image)
	})
}

var playerResponse = `{
	"data": {
		"player_id": 219591,
		"team_id": 1,
		"country_id": 1190,
		"position_id": 2,
		"common_name": "F. Balbuena",
		"fullname": "F. Balbuena",
		"firstname": "Fabián Cornelio",
		"lastname": "Balbuena González",
		"nationality": "Paraguay",
		"birthdate": "23\/08\/1991",
		"birthcountry": "Paraguay",
		"birthplace": "Ciudad del Este",
		"height": "188 cm",
		"weight": "82 kg",
		"image_path": "https://cdn.sportmonks.com/images/soccer/player/1/1.png"
	}
}`
