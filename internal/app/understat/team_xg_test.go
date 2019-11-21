package understat_test

import (
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app/understat"
	parser "github.com/statistico/statistico-understat-parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixtureTeamXGRequester_FixtureTeamXGByLeagueAndSeason(t *testing.T) {
	t.Run("returns a channel containing understat fixture struct", func(t *testing.T) {
		p := parser.Parser{BaseURL:"https://understat.com"}
		logger, _ := test.NewNullLogger()

		requester := understat.NewFixtureTeamXGRequester(&p, logger)

		fixtures := requester.FixtureTeamXGByLeagueAndSeason("EPL", "2019")

		_ = <- fixtures

		assert.Equal(t, 380, len(fixtures))
	})
}
