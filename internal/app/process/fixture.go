package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const fixture = "fixture"
const fixtureCurrentSeason = "fixture:current-season"

// FixtureProcessor fetches data from external data source using the FixtureRequester
// before persisting to the storage engine using the FixtureRepository
type FixtureProcessor struct {
	fixtureRepo app.FixtureRepository
	seasonRepo  app.SeasonRepository
	requester   app.FixtureRequester
	logger      *logrus.Logger
}

func (f FixtureProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case fixture:
		go f.processAllSeasons(done)
	case fixtureCurrentSeason:
		go f.processCurrentSeason(done)
	default:
		f.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (f FixtureProcessor) processAllSeasons(done chan bool) {
	ids, err := f.seasonRepo.IDs()

	if err != nil {
		f.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := f.requester.FixturesBySeasonIDs(ids)

	go f.persistFixtures(ch, done)
}

func (f FixtureProcessor) processCurrentSeason(done chan bool) {
	ids, err := f.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		f.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := f.requester.FixturesBySeasonIDs(ids)

	go f.persistFixtures(ch, done)
}

func (f FixtureProcessor) persistFixtures(ch <-chan *app.Fixture, done chan bool) {
	for fixture := range ch {
		f.persist(fixture)
	}

	done <- true
}

func (f FixtureProcessor) persist(x *app.Fixture) {
	_, err := f.fixtureRepo.ByID(x.ID)

	if err != nil {
		if err := f.fixtureRepo.Insert(x); err != nil {
			f.logger.Warningf("Error '%s' occurred when inserting fixture struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := f.fixtureRepo.Update(x); err != nil {
		f.logger.Warningf("Error '%s' occurred when updating fixture struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewFixtureProcessor(f app.FixtureRepository, s app.SeasonRepository, r app.FixtureRequester, log *logrus.Logger) *FixtureProcessor {
	return &FixtureProcessor{fixtureRepo: f, seasonRepo: s, requester: r, logger: log}
}
