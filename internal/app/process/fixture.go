package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
)

const fixturesCurrentSeason = "fixtures:current-season"
const fixturesBySeasonId = "fixtures:by-season-id"
const fixturesByCompetitionId = "fixtures:by-competition-id"

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
	case fixturesCurrentSeason:
		go f.processCurrentSeason(done)
	case fixturesBySeasonId:
		id, _ := strconv.Atoi(option)
		go f.processSeason(uint64(id), done)
	case fixturesByCompetitionId:
		id, _ := strconv.Atoi(option)
		go f.processSeason(uint64(id), done)
	default:
		f.logger.Fatalf("Command %s is not supported", command)
		return
	}
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

func (f FixtureProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := f.requester.FixturesBySeasonIDs([]uint64{seasonID})

	go f.persistFixtures(ch, done)
}

func (f FixtureProcessor) processCompetition(competitionID uint64, done chan bool) {
	seasons, err := f.seasonRepo.ByCompetitionId(competitionID, "name_asc")

	if err != nil {
		f.logger.Fatalf("Error when retrieving seasons in fixture processor: %s", err.Error())
		return
	}

	var ids []uint64

	for _, season := range seasons {
		ids = append(ids, season.ID)
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
	if x.Status != nil && (*x.Status == "Deleted" || *x.Status == "POSTP"){
		if err := f.fixtureRepo.Delete(x.ID); err != nil {
			f.logger.Warningf("Error '%s' occurred when delete fixture: %d\n,", err.Error(), x.ID)
		}

		return
	}

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
