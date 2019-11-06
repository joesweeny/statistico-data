package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"strings"
	"time"
)

const resultById = "result:by-id"
const resultBySeasonId = "result:by-season-id"
const resultToday = "result:today"

type ResultProcessor struct {
	resultRepo app.ResultRepository
	fixtureRepo app.FixtureRepository
	requester app.ResultRequester
	clock clockwork.Clock
	logger     *logrus.Logger
}

func (r ResultProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case resultById:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go r.processByID(done, uint64(id))
		}
	case resultBySeasonId:
		id, _ := strconv.Atoi(option)
		go r.processSeason(done, uint64(id))
	case resultToday:
		go r.processTodayResults(done)
	default:
		r.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (r ResultProcessor) processByID(done chan bool, id uint64) {
	fix, err := r.fixtureRepo.ByID(id)

	if err != nil {
		r.logger.Fatalf("Error when retrieving fixtures for ID: %d, %s", id, err.Error())
		return
	}

	ch := r.requester.ResultsByFixtureIDs([]uint64{fix.ID})

	go r.persistResults(ch, done)
}

func (r ResultProcessor) processSeason(done chan bool, seasonID uint64) {
	fix, err := r.fixtureRepo.BySeasonID(seasonID)

	if err != nil {
		r.logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", seasonID, err.Error())
		return
	}

	var ids []uint64

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	ch := r.requester.ResultsByFixtureIDs(ids)

	go r.persistResults(ch, done)
}

func (r ResultProcessor) processTodayResults(done chan bool) {
	now := r.clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := r.fixtureRepo.IDsBetween(from, to)

	if err != nil {
		r.logger.Fatalf("Error when retrieving fixture ids in event processor: %s", err.Error())
		return
	}

	ch := r.requester.ResultsByFixtureIDs(ids)

	go r.persistResults(ch, done)
}

func (r ResultProcessor) persistResults(ch <-chan *app.Result, done chan bool) {
	for result := range ch {
		r.persist(result)
	}

	done <- true
}

func (r ResultProcessor) persist(x *app.Result) {
	_, err := r.resultRepo.ByFixtureID(x.FixtureID)

	if err != nil {
		if err := r.resultRepo.Insert(x); err != nil {
			r.logger.Warningf("Error '%s' occurred when inserting result struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := r.resultRepo.Update(x); err != nil {
		r.logger.Warningf("Error '%s' occurred when updating result struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewResultProcessor(r app.ResultRepository, f app.FixtureRepository, q app.ResultRequester, c clockwork.Clock, log *logrus.Logger) *ResultProcessor {
	return &ResultProcessor{resultRepo: r, fixtureRepo: f, requester: q, clock: c, logger: log}
}
