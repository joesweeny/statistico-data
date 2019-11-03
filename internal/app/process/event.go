package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/fixture"
	"sync"
	"time"
)

// EventProcessor fetches data from external data source using the EventRequester
// before persisting to the storage engine using the EventRepository.
type EventProcessor struct {
	eventRepo app.EventRepository
	fixtureRepo fixture.Repository
	requester app.EventRequester
	clock clockwork.Clock
	logger     *logrus.Logger
}

func (e EventProcessor) Process(command string, option string, done chan bool) {

}

func (e EventProcessor) processTodayEvents(g <-chan *app.GoalEvent, s <-chan *app.SubstitutionEvent, done chan bool) {
	now := e.clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := e.fixtureRepo.IdsBetween(from, to)

	if err != nil {
		e.logger.Fatalf("Error when retrieving fixture ids in event processor: %s", err.Error())
		return
	}


}

func (e EventProcessor) persistEvents(g <-chan *app.GoalEvent, s <-chan *app.SubstitutionEvent, done chan bool) {
	wg := sync.WaitGroup{}

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
