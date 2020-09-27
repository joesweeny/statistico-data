package process

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	understat "github.com/statistico/statistico-understat-parser"
	"strconv"
)

const fixtureXG = "fixture-xg"
const fixtureXGCurrentSeason = "fixture-xg:current-season"

var currentSeason = map[string]map[int]string {
	"Bundesliga": {
		17361: "2020",
	},
	"EPL": {
		17420: "2020",
	},
	"La liga": {
		17480: "2020",
	},
	"Ligue_1": {
		17160: "2020",
	},
	"Serie A": {
		17488: "2020",
	},
}

var historicSeasons =  map[string]map[int]string {
	"Bundesliga": {
		16264: "2019",
		13005: "2018",
		8026: "2017",
		219: "2016",
		218: "2015",
		217: "2014",
	},
	"EPL": {
		16036: "2019",
		12962: "2018",
		6397: "2017",
		13: "2016",
		10: "2015",
	},
	"La liga": {
		16326: "2019",
		13133: "2018",
		8442: "2017",
		853: "2016",
		2063: "2015",
		2061: "2014",
	},
	"Ligue_1": {
		16043: "2019",
		12935: "2018",
		6405: "2017",
		765: "2016",
		1390: "2015",
		1389: "2014",
	},
	"Serie A": {
		16415: "2019",
		13158: "2018",
		8557: "2017",
		802: "2016",
		1584: "2015",
		1583: "2014",
	},
}

var teamMapper = map[string]string {
	"GFC Ajaccio": "Gazélec Ajaccio",
	"AC Milan": "Milan",
	"Alaves": "Deportivo Alavés",
	"Atletico Madrid": "Atlético Madrid",
	"Almeria": "Almería",
	"Bayern Munich": "Bayern München",
	"Cadiz": "Cádiz",
	"Celta Vigo": "Celta de Vigo",
	"Cordoba": "Córdoba",
	"Deportivo La Coruna": "Deportivo La Coruña",
	"FC Cologne": "Köln",
	"Bournemouth": "AFC Bournemouth",
	"Borussia M.Gladbach": "Borussia M'gladbach",
	"Eibar": "SD Eibar",
	"Evian Thonon Gaillard": "Evian TG",
	"Fortuna Duesseldorf": "Fortuna Düsseldorf",
	"Hertha Berlin": "Hertha BSC",
	"Leganes": "Leganés",
	"Lyon": "Olympique Lyonnais",
	"Malaga": "Málaga",
	"Marseille": "Olympique Marseille",
	"Nimes": "Nîmes",
	"Nuernberg": "Nürnberg",
	"Parma Calcio 1913": "Parma",
	"RasenBallsport Leipzig": "RB Leipzig",
	"Saint-Etienne": "Saint-Étienne",
	"SC Bastia": "Bastia",
	"SD Huesca": "Huesca",
	"SPAL 2013": "SPAL",
	"Sporting Gijon": "Sporting Gijón",
	"Verona": "Hellas Verona",
	"VfB Stuttgart": "Stuttgart",
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

func (f FixtureTeamXGProcessor) processFixtures(done chan bool, seasons map[string]map[int]string) {
	for k, v := range seasons {
		for id, year := range v {
			fix, err := f.parser.LeagueFixtures(k, year)

			if err != nil {
				f.logger.Warnf("error fetching league xg data. League %s, Season %s", k, year)
				continue
			}

			f.parseFixtures(fix, uint64(id))
		}
	}

	done <- true
}

func (f FixtureTeamXGProcessor) parseFixtures(fixtures []understat.Fixture, seasonID uint64) {
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

		f.createNew(fix, seasonID)
	}
}

func (f FixtureTeamXGProcessor) createNew(u understat.Fixture, seasonID uint64) {
	fixture, err := f.parseFixture(u, seasonID)

	if err != nil {
		f.logger.Warnf("error creating new fixture team xg struct. %s", err.Error())
		return
	}

	id, err := strconv.Atoi(u.ID)

	if err != nil {
		f.logger.Warnf("error parsing string to int creating new fixture team xg struct. %s", err.Error())
		return
	}

	home, err1 := parseFloat(u.XG.Home)
	away, err2 := parseFloat(u.XG.Away)

	if err1 != nil || err2 != nil {
		f.logger.Warnf("unable to parse float when processing fixture id %d", fixture.ID)
		return
	}

	xg := &app.FixtureTeamXG{
		ID:        uint64(id),
		FixtureID: fixture.ID,
		Home:      home,
		Away:      away,
	}

	if err := f.xGRepo.Insert(xg); err != nil {
		f.logger.Warnf("error inserting fixture team xg %s, fixture id %d", u.ID, xg.FixtureID)
	}
}

func (f FixtureTeamXGProcessor) updateExisting(xg *app.FixtureTeamXG, u understat.Fixture) {
	home, err1 := parseFloat(u.XG.Home)
	away, err2 := parseFloat(u.XG.Away)

	if err1 != nil || err2 != nil {
		f.logger.Warnf("unable to parse float when processing fixture id %d", xg.FixtureID)
		return
	}

	xg.Home = home
	xg.Away = away

	if err := f.xGRepo.Update(xg); err != nil {
		f.logger.Warnf("error update fixture team xg %d, fixture id %d", xg.ID, xg.FixtureID)
	}
}

func (f FixtureTeamXGProcessor) parseFixture(u understat.Fixture, seasonID uint64) (*app.Fixture, error) {
	home := parseTeam(u.Home.Title)
	away := parseTeam(u.Away.Title)

	query := app.FixtureRepositoryQuery{
		HomeTeamNameLike: &home,
		AwayTeamNameLike: &away,
		SeasonID: &seasonID,
	}

	fixs, err := f.fixtureRepo.Get(query)

	if err != nil || len(fixs) == 0 {
		return nil, fmt.Errorf("unable to find matching fixture xg for understat ID %s, home %s, away %s", u.ID, home, away)
	}

	return &fixs[0], nil
}

func parseFloat(str *string) (*float32, error) {
	if str == nil {
		return nil, nil
	}

	float, err := strconv.ParseFloat(*str, 32)

	if err != nil {
		return nil, err
	}

	f := float32(float)

	return &f, nil
}

func parseTeam(team string) string {
	value, ok := teamMapper[team]

	if ok {
		team = value
	}

	return team
}

func NewFixtureTeamXGProcessor(
	r app.FixtureTeamXGRepository,
	f app.FixtureRepository,
	p *understat.Parser,
	l *logrus.Logger,
) *FixtureTeamXGProcessor {
	return &FixtureTeamXGProcessor{xGRepo: r, fixtureRepo: f, parser: p, logger: l}
}
