package event

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/fixture"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const eventsByResultId = "events:by-result-id"
const eventsBySeasonId = "events:by-season-id"
const eventsToday = "events:today"
const callLimit = 1500

var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	Repository
	Factory
	Logger 		*log.Logger
	FixtureRepo fixture.Repository
	Client      *sportmonks.Client
}

func (p Processor) Process(command string, option string, done chan bool) {
	switch command {
	case eventsByResultId:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go p.byId(done, id)
		}
	case eventsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.bySeasonId(done, id)
	case eventsToday:
		go p.statsToday(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) byId(done chan bool, id int) {
	fix, err := p.FixtureRepo.ById(uint64(id))

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []int{fix.ID}

	go p.processEvents(ids, done)
}

func (p Processor) bySeasonId(done chan bool, id int) {
	// Adding a Clock.Now() here is a bit hacky. Redo by dynamically handling this
	fix, err := p.FixtureRepo.BySeasonId(int64(id), p.Clock.Now())

	if err != nil {
		p.Logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", id, err.Error())
		return
	}

	var ids []int

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	go p.processEvents(ids, done)
}

func (p Processor) statsToday(done chan bool) {
	now := p.Clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := p.FixtureRepo.IdsBetween(from, to)

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processEvents(ids, done)
}

func (p Processor) processEvents(ids []int, done chan bool) {
	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseEvents(results, done)
}

func (p Processor) callClient(ids []int, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"goals,substitutions"}

	for _, id := range ids {
		if *c >= callLimit {
			p.Logger.Printf("Api call limited reached %d calls\n", *c)
			break
		}

		res, err := p.Client.FixtureById(id, q, 5)

		*c++

		if err != nil {
			p.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		ch <- res.Data
	}

	close(ch)
}

func (p Processor) parseEvents(ch <-chan sportmonks.Fixture, done chan bool) {
	for x := range ch {
		p.handleEvents(x)
	}

	waitGroup.Wait()

	done <- true
}

func (p Processor) handleEvents(fix sportmonks.Fixture) {
	for _, sub := range fix.Subs.Data {
		waitGroup.Add(1)

		go func(e sportmonks.SubstitutionEvent) {
			p.processSubstitutionEvent(&e)
			defer waitGroup.Done()
		}(sub)
	}

	for _, goal := range fix.Goals.Data {
		waitGroup.Add(1)

		go func(e sportmonks.GoalEvent) {
			p.processGoalEvent(&e)
			defer waitGroup.Done()
		}(goal)
	}
}

func (p Processor) processGoalEvent(s *sportmonks.GoalEvent) {
	if _, err := p.Repository.GoalEventById(s.ID); err != ErrNotFound {
		return
	}

	created := p.Factory.createGoalEvent(s)

	if err := p.InsertGoalEvent(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Goal Event struct: %+v\n,", err.Error(), created)
	}

	return
}

func (p Processor) processSubstitutionEvent(s *sportmonks.SubstitutionEvent) {
	if _, err := p.Repository.SubstitutionEventById(s.ID); err != ErrNotFound {
		return
	}

	created := p.Factory.createSubstitutionEvent(s)

	if err := p.InsertSubstitutionEvent(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Substitution Event struct: %+v\n,", err.Error(), created)
	}

	return
}
