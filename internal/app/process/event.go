package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"sync"
	"time"
)

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
