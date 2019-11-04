package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"sync"
	"time"
)

const eventsByResultId = "events:by-result-id"
const eventsBySeasonId = "events:by-season-id"
const eventsToday = "events:today"

// EventProcessor fetches data from external data source using the EventRequester
// before persisting to the storage engine using the EventRepository.
type EventProcessor struct {
	eventRepo app.EventRepository
	fixtureRepo app.FixtureRepository
	requester app.EventRequester
	clock clockwork.Clock
	logger     *logrus.Logger
}

func (e EventProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case eventsByResultId:
		id, _ := strconv.Atoi(option)
		go e.processEventsByFixtureID(uint64(id), done)
	case eventsBySeasonId:
		id, _ := strconv.Atoi(option)
		go e.processEventsBySeasonID(uint64(id), done)
	case eventsToday:
		go e.processTodayEvents(done)
	default:
		e.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (e EventProcessor) processEventsByFixtureID(id uint64, done chan bool) {
	fix, err := e.fixtureRepo.ByID(id)

	if err != nil {
		e.logger.Fatalf("Error when retrieving fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []uint64{fix.ID}

	goals, subs := e.requester.EventsByFixtureIDs(ids)

	go e.parseEvents(goals, subs, done)
}

func (e EventProcessor) processEventsBySeasonID(id uint64, done chan bool) {
	fix, err := e.fixtureRepo.BySeasonID(id)

	if err != nil {
		e.logger.Fatalf("Error when retrieving fixtures for season ID: %d, %s", id, err.Error())
		return
	}

	var ids []uint64

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	goals, subs := e.requester.EventsByFixtureIDs(ids)

	go e.parseEvents(goals, subs, done)
}

func (e EventProcessor) processTodayEvents(done chan bool) {
	now := e.clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := e.fixtureRepo.IDsBetween(from, to)

	if err != nil {
		e.logger.Fatalf("Error when retrieving fixture ids in event processor: %s", err.Error())
		return
	}

	goals, subs := e.requester.EventsByFixtureIDs(ids)

	go e.parseEvents(goals, subs, done)
}

func (e EventProcessor) parseEvents(g <-chan *app.GoalEvent, s <-chan *app.SubstitutionEvent, done chan bool) {
	var wg = sync.WaitGroup{}

	wg.Add(2)

	go func(g <-chan *app.GoalEvent) {
		for goal := range g {
			e.persistGoalEvent(goal)
		}

		wg.Done()
	}(g)

	go func(g <-chan *app.SubstitutionEvent) {
		for sub := range s {
			e.persistSubstitutionEvent(sub)
		}

		wg.Done()
	}(s)

	done <- true
}

func (e EventProcessor) persistGoalEvent(x *app.GoalEvent) {
	if _, err := e.eventRepo.GoalEventByID(x.ID); err != nil {
		return
	}

	if err := e.eventRepo.InsertGoalEvent(x); err != nil {
		e.logger.Warningf("Error '%s' occurred when inserting goal event struct: %+v\n,", err.Error(), *x)
	}
}

func (e EventProcessor) persistSubstitutionEvent(x *app.SubstitutionEvent) {
	if _, err := e.eventRepo.SubstitutionEventByID(x.ID); err != nil {
		return
	}

	if err := e.eventRepo.InsertSubstitutionEvent(x); err != nil {
		e.logger.Warningf("Error '%s' occurred when inserting substitution event struct: %+v\n,", err.Error(), *x)
	}
}
