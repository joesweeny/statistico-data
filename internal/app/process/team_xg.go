package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	understat "github.com/statistico/statistico-understat-parser"
)

const fixtureXG = "fixture-xg"
const fixtureXGCurrentSeason = "fixture-xg:current-season"

var currentSeason map[string]string{
	"EPL": "2019"
}

var historicSeasons map[string]string{
	"EPL": "2018",
	"EPL": "2017",
	"EPL": "2016",
	"EPL": "2015",
	"EPL": "2014",
}

type FixtureTeamXGProcessor struct {
	xGRepo app.FixtureTeamXGRepository
	fixtureRepo app.FixtureRepository
	parser *understat.Parser
	logger *logrus.Logger
}

func (f FixtureTeamXGProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case fixtureXG:
		f.processFixtures(done, historicSeasons)
	case fixtureXGCurrentSeason:
		f.processFixtures(done, currentSeason)
	default:
		f.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (f FixtureTeamXGProcessor) processFixtures(done chan bool, seasons map[string]string) {
	for k, v := range seasons {
		fix, err := f.parser.LeagueFixtures(k, v)

		if err != nil {
			// Do something
			continue
		}

		f.parseFixtures(fix)
	}

	done <- true
}

func (f FixtureTeamXGProcessor) parseFixtures(fix []understat.Fixture) {
	// Check if ID assists in XG Repo

	// If exists then update struct then update in repo

	// If not exists then check fixture repo for matching fixture
	   // - Error if not found
	   // - or hydrate new struct and insert
}
