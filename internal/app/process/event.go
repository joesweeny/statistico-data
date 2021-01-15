package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"sync"
)

const events = "events"
const eventsCurrentSeason = "events:current-season"
const eventsBySeasonId = "events:by-season-id"


// EventProcessor fetches data from external data source using the EventRequester
// before persisting to the storage engine using the EventRepository.
type EventProcessor struct {
	eventRepo   app.EventRepository
	seasonRepo  app.SeasonRepository
	requester   app.EventRequester
	clock       clockwork.Clock
	logger      *logrus.Logger
}

func (e EventProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case events:
		go e.processAllSeasons(done)
	case eventsCurrentSeason:
		go e.processCurrentSeason(done)
	case eventsBySeasonId:
		id, _ := strconv.Atoi(option)
		go e.processEventsBySeasonID(uint64(id), done)
	default:
		e.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (e EventProcessor) processAllSeasons(done chan bool) {
	ids, err := e.seasonRepo.IDs()

	if err != nil {
		e.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	goals, subs, cards := e.requester.EventsBySeasonIDs(ids)

	go e.parseEvents(goals, subs, cards, done)
}

func (e EventProcessor) processCurrentSeason(done chan bool) {
	ids, err := e.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		e.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	goals, subs, cards := e.requester.EventsBySeasonIDs(ids)

	go e.parseEvents(goals, subs, cards, done)
}

func (e EventProcessor) processEventsBySeasonID(seasonID uint64, done chan bool) {
	goals, subs, cards := e.requester.EventsBySeasonIDs([]uint64{seasonID})

	go e.parseEvents(goals, subs, cards, done)
}

func (e EventProcessor) parseEvents(g <-chan *app.GoalEvent, s <-chan *app.SubstitutionEvent, c <-chan *app.CardEvent, done chan bool) {
	var wg = sync.WaitGroup{}

	wg.Add(3)

	go func(c <-chan *app.CardEvent) {
		for card := range c {
			e.persistCardEvent(card)
		}

		wg.Done()
	}(c)

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

	wg.Wait()

	done <- true
}

func (e EventProcessor) persistCardEvent(x *app.CardEvent) {
	if _, err := e.eventRepo.CardEventByID(x.ID); err == nil {
		return
	}

	if err := e.eventRepo.InsertCardEvent(x); err != nil {
		e.logger.Warningf("Error '%s' occurred when inserting card event struct: %+v\n,", err.Error(), *x)
	}
}

func (e EventProcessor) persistGoalEvent(x *app.GoalEvent) {
	if _, err := e.eventRepo.GoalEventByID(x.ID); err == nil {
		return
	}

	if err := e.eventRepo.InsertGoalEvent(x); err != nil {
		e.logger.Warningf("Error '%s' occurred when inserting goal event struct: %+v\n,", err.Error(), *x)
	}
}

func (e EventProcessor) persistSubstitutionEvent(x *app.SubstitutionEvent) {
	if _, err := e.eventRepo.SubstitutionEventByID(x.ID); err == nil {
		return
	}

	if err := e.eventRepo.InsertSubstitutionEvent(x); err != nil {
		e.logger.Warningf("Error '%s' occurred when inserting substitution event struct: %+v\n,", err.Error(), *x)
	}
}

func NewEventProcessor(r app.EventRepository, s app.SeasonRepository, q app.EventRequester, c clockwork.Clock, log *logrus.Logger) *EventProcessor {
	return &EventProcessor{eventRepo: r, seasonRepo: s, requester: q, clock: c, logger: log}
}
