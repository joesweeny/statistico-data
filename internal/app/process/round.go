package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const round = "round"
const roundCurrentSeason = "round:current-season"

// Process fetches data from external data source using the RoundRequester
// before persisting to the storage engine using the RoundRepository
type RoundProcessor struct {
	roundRepo app.RoundRepository
	seasonRepo app.SeasonRepository
	requester app.RoundRequester
	logger     *logrus.Logger
}

func (r RoundProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case round:
		go r.processAllSeasons(done)
	case roundCurrentSeason:
		go r.processCurrentSeason(done)
	default:
		r.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (r RoundProcessor) processAllSeasons(done chan bool) {
	ids, err := r.seasonRepo.IDs()

	if err != nil {
		r.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := r.requester.RoundsBySeasonIDs(ids)

	go r.persistRounds(ch, done)
}

func (r RoundProcessor) processCurrentSeason(done chan bool) {
	ids, err := r.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		r.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := r.requester.RoundsBySeasonIDs(ids)

	go r.persistRounds(ch, done)
}

func (r RoundProcessor) persistRounds(ch <-chan *app.Round, done chan bool) {
	for round := range ch {
		r.persist(round)
	}

	done <- true
}

func (r RoundProcessor) persist(x *app.Round) {
	_, err := r.roundRepo.ByID(x.ID)

	if err != nil {
		if err := r.roundRepo.Insert(x); err != nil {
			r.logger.Warningf("Error '%s' occurred when inserting round struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := r.roundRepo.Update(x); err != nil {
		r.logger.Warningf("Error '%s' occurred when updating round struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewRoundProcessor(r app.RoundRepository, s app.SeasonRepository, q app.RoundRequester, log *logrus.Logger) *RoundProcessor {
	return &RoundProcessor{roundRepo: r, seasonRepo: s, requester: q, logger: log}
}
