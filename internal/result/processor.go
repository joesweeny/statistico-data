package result

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/fixture"
	"log"
	"strconv"
	"strings"
	"time"
)

const result = "result"
const resultById = "result:by-id"
const resultBySeasonId = "result:by-season-id"
const resultToday = "result:today"
const callLimit = 1500

var counter int

type Processor struct {
	Repository
	FixtureRepo fixture.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
	Clock  clockwork.Clock
}

func (p Processor) Process(command string, option string, done chan bool) {
	switch command {
	case result:
		go p.allResults(done)
	case resultById:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go p.byId(done, id)
		}
	case resultBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.bySeasonId(done, id)
	case resultToday:
		go p.resultsToday(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) allResults(done chan bool) {
	ids, err := p.FixtureRepo.Ids()

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processResults(ids, done)
}

func (p Processor) byId(done chan bool, id int) {
	fix, err := p.FixtureRepo.ById(uint64(id))

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []int{fix.ID}

	go p.processResults(ids, done)
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

	go p.processResults(ids, done)
}

func (p Processor) resultsToday(done chan bool) {
	now := p.Clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := p.FixtureRepo.IdsBetween(from, to)

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processResults(ids, done)
}

func (p Processor) processResults(ids []int, done chan bool) {
	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseResults(results, done)
}

func (p Processor) callClient(ids []int, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	var q []string

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

func (p Processor) parseResults(ch <-chan sportmonks.Fixture, done chan bool) {
	for x := range ch {
		p.handleResult(x)
	}

	done <- true
}

func (p Processor) handleResult(fix sportmonks.Fixture) {
	x, err := p.Repository.GetByFixtureId(fix.ID)

	if err == ErrNotFound {
		created := p.Factory.createResult(&fix)

		if err := p.Repository.Insert(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Result struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := p.Factory.updateResult(&fix, x)

	if err := p.Repository.Update(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Result struct: %+v\n,", err.Error(), updated)
	}

	return
}
