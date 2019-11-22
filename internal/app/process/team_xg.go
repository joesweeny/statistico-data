package process

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	understat "github.com/statistico/statistico-understat-parser"
	"strconv"
	"time"
)

const fixtureXG = "fixture-xg"
const fixtureXGCurrentSeason = "fixture-xg:current-season"

var currentSeason = map[string][]string {
	"2019": {
		"EPL",
	},
}

var historicSeasons =  map[string][]string {
	"2018": {
		"EPL",
	},
	"2017": {
		"EPL",
	},
	"2016": {
		"EPL",
	},
	"2015": {
		"EPL",
	},
	"2014": {
		"EPL",
	},
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

func (f FixtureTeamXGProcessor) processFixtures(done chan bool, seasons map[string][]string) {
	for k, v := range seasons {
		for _, league := range v {
			fix, err := f.parser.LeagueFixtures(league, k)

			if err != nil {
				f.logger.Warnf("error fetching league xg data. League %s, Season %s", league, k)
				continue
			}

			f.parseFixtures(fix)
		}
	}

	done <- true
}

func (f FixtureTeamXGProcessor) parseFixtures(fixtures []understat.Fixture) {
	for _, fix := range fixtures {
		id, err := strconv.Atoi(fix.ID)

		if err != nil {
			f.logger.Fatalf("error parsing string to int in FixtureTeamXGProcessor. %s", fix.ID)
		}

		xg, err := f.xGRepo.ByID(uint64(id))

		if err == nil {
			f.updateExisting(xg, fix)
			continue
		}

		f.createNew(fix)
	}
}

func (f FixtureTeamXGProcessor) createNew(u understat.Fixture) {
	fixture, err := f.parseFixture(u)

	if err != nil {
		f.logger.Warnf("error creating new fixture team xg struct. %s", err.Error())
		return
	}

	id, err := strconv.Atoi(u.ID)

	if err != nil {
		f.logger.Warnf("error parsing string to int creating new fixture team xg struct. %s", err.Error())
		return
	}

	home, err1 := parseFloat(*u.XG.Home)
	away, err2 := parseFloat(*u.XG.Away)

	if err1 != nil || err2 != nil {
		f.logger.Warnf("unable to parse float when processing fixture id %d", fixture.ID)
		return
	}

	xg := &app.FixtureTeamXG{
		ID:        uint64(id),
		FixtureID: fixture.ID,
		Home:      &home,
		Away:      &away,
	}

	if err := f.xGRepo.Insert(xg); err != nil {
		f.logger.Warnf("error inserting fixture team xg %s, fixture id %d", u.ID, xg.FixtureID)
	}
}

func (f FixtureTeamXGProcessor) updateExisting(xg *app.FixtureTeamXG, u understat.Fixture) {
	home, err1 := parseFloat(*u.XG.Home)
	away, err2 := parseFloat(*u.XG.Away)

	if err1 != nil || err2 != nil {
		f.logger.Warnf("unable to parse float when processing fixture id %d", xg.FixtureID)
		return
	}

	xg.Home = &home
	xg.Away = &away

	if err := f.xGRepo.Update(xg); err != nil {
		f.logger.Warnf("error update fixture team xg %d, fixture id %d", xg.ID, xg.FixtureID)
	}
}

func (f FixtureTeamXGProcessor) parseFixture(u understat.Fixture) (*app.Fixture, error) {
	home := u.Home.Title[0:5]
	away := u.Away.Title[0:5]

	from, err1 := parseDateTime(u.DateTime, -2*time.Hour)
	to, err2 := parseDateTime(u.DateTime, 2*time.Hour)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("unable to parse date when processing understat id %s", u.ID)
	}

	query := app.FixtureRepositoryQuery{
		HomeTeamNameLike: &home,
		AwayTeamNameLike: &away,
		DateFrom: from,
		DateTo: to,
	}

	fixs, err := f.fixtureRepo.Get(query)

	if err != nil || len(fixs) == 0 {
		return nil, fmt.Errorf("unable to find matching fixture xg for understat ID %s", u.ID)
	}

	return &fixs[0], nil
}

func parseFloat(str string) (float32, error) {
	float, err := strconv.ParseFloat(str, 32)

	if err != nil {
		return 0, err
	}

	return float32(float), nil
}

func parseDateTime(d string, duration time.Duration) (*time.Time, error) {
	layout := "2006-01-02 15:04:05"

	date, err := time.Parse(layout, d)

	if err != nil {
		return nil, err
	}

	date = date.Add(duration)

	return &date, nil
}

func NewFixtureTeamXGProcessor(r app.FixtureTeamXGRepository, f app.FixtureRepository, p *understat.Parser, l *logrus.Logger) *FixtureTeamXGProcessor {
	return &FixtureTeamXGProcessor{xGRepo: r, fixtureRepo: f, parser: p, logger: l}
}
