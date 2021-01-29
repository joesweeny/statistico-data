package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
)

const resultsCurrentSeason = "results:current-season"
const resultsBySeasonId = "results:by-season-id"
const resultsByCompetitionId = "results:by-competition-id"

type ResultProcessor struct {
	resultRepo  app.ResultRepository
	seasonRepo  app.SeasonRepository
	requester   app.ResultRequester
	clock       clockwork.Clock
	logger      *logrus.Logger
}

func (r ResultProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case resultsCurrentSeason:
		go r.processCurrentSeason(done)
	case resultsBySeasonId:
		id, _ := strconv.Atoi(option)
		go r.processSeason(uint64(id), done)
	case resultsByCompetitionId:
		id, _ := strconv.Atoi(option)
		go r.processCompetition(uint64(id), done)
	default:
		r.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (r ResultProcessor) processCurrentSeason(done chan bool) {
	ids, err := r.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		r.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := r.requester.ResultsBySeasonIDs(ids)

	go r.persistResults(ch, done)
}

func (r ResultProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := r.requester.ResultsBySeasonIDs([]uint64{seasonID})

	go r.persistResults(ch, done)
}

func (r ResultProcessor) processCompetition(competitionID uint64, done chan bool) {
	seasons, err := r.seasonRepo.ByCompetitionId(competitionID, "name_asc")

	if err != nil {
		r.logger.Fatalf("Error when retrieving seasons in result processor: %s", err.Error())
		return
	}

	var ids []uint64

	for _, season := range seasons {
		ids = append(ids, season.ID)
	}

	ch := r.requester.ResultsBySeasonIDs(ids)

	go r.persistResults(ch, done)
}

func (r ResultProcessor) persistResults(ch <-chan app.Result, done chan bool) {
	for result := range ch {
		r.persist(result)
	}

	done <- true
}

func (r ResultProcessor) persist(x app.Result) {
	_, err := r.resultRepo.ByFixtureID(x.FixtureID)

	if err != nil {
		if err := r.resultRepo.Insert(&x); err != nil {
			r.logger.Errorf("Error '%s' occurred when inserting result struct: %+v\n,", err.Error(), x)
		}

		return
	}

	if err := r.resultRepo.Update(&x); err != nil {
		r.logger.Errorf("Error '%s' occurred when updating result struct: %+v\n,", err.Error(), x)
	}

	return
}

func NewResultProcessor(r app.ResultRepository, f app.SeasonRepository, q app.ResultRequester, c clockwork.Clock, log *logrus.Logger) *ResultProcessor {
	return &ResultProcessor{resultRepo: r, seasonRepo: f, requester: q, clock: c, logger: log}
}
