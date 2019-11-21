package understat

import (
	"github.com/sirupsen/logrus"
	parser "github.com/statistico/statistico-understat-parser"
)

type FixtureTeamXGRequester struct {
	parser 	*parser.Parser
	logger  *logrus.Logger
}

func (f FixtureTeamXGRequester) FixtureTeamXGByLeagueAndSeason(league, season string) <-chan *parser.Fixture {
	ch := make(chan *parser.Fixture, 380)

	go f.parseFixtures(league, season, ch)

	return ch
}

func (f FixtureTeamXGRequester) parseFixtures(league, season string, ch chan<- *parser.Fixture) {
	defer close(ch)

	fixtures, err := f.parser.LeagueFixtures(league, season)

	if err != nil {
		panic(err)
		f.logger.Fatalf("Error parse fixture team xg from understat '%s'", err.Error())
		return
	}

	for _, fix := range fixtures {
		ch <- &fix
	}
}

func NewFixtureTeamXGRequester(p *parser.Parser, logger *logrus.Logger) *FixtureTeamXGRequester {
	return &FixtureTeamXGRequester{parser: p, logger: logger}
}
