package result

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statistico-data/internal/event"
	"github.com/joesweeny/statistico-data/internal/fixture"
	"github.com/joesweeny/statistico-data/internal/stats"
	"log"
	"sync"
	"time"
)

const result = "result"
const resultToday = "result:today"
const callLimit = 1500

var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	Repository
	FixtureRepo fixture.Repository
	Factory
	Client          *sportmonks.Client
	Logger          *log.Logger
	PlayerProcessor stats.PlayerProcessor
	TeamProcessor   stats.TeamProcessor
	EventProcessor  event.Processor
}

func (p Processor) Process(command string, done chan bool) {
	switch command {
	case result:
		go p.allResults(done)
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
	q := []string{"lineup,bench,stats,goals,substitutions"}

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

	waitGroup.Wait()

	done <- true
}

func (p Processor) handleResult(fix sportmonks.Fixture) {
	waitGroup.Add(5)

	go func(stats []sportmonks.TeamStats) {
		p.handleTeams(stats)
		defer waitGroup.Done()
	}(fix.TeamStats.Data)

	go func(players []sportmonks.LineupPlayer) {
		p.handlePlayers(players, false)
		defer waitGroup.Done()
	}(fix.Lineup.Data)

	go func(players []sportmonks.LineupPlayer) {
		p.handlePlayers(players, true)
		defer waitGroup.Done()
	}(fix.Bench.Data)

	go func(goals []sportmonks.GoalEvent) {
		p.handleGoalEvents(goals)
		defer waitGroup.Done()
	}(fix.Goals.Data)

	go func(subs []sportmonks.SubstitutionEvent) {
		p.handleSubstitutionEvents(subs)
		defer waitGroup.Done()
	}(fix.Subs.Data)

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

func (p Processor) handleTeams(t []sportmonks.TeamStats) {
	for _, team := range t {
		waitGroup.Add(1)

		go func(stats sportmonks.TeamStats) {
			p.TeamProcessor.ProcessTeamStats(&stats)
			defer waitGroup.Done()
		}(team)
	}
}

func (p Processor) handlePlayers(lineups []sportmonks.LineupPlayer, bench bool) {
	for _, player := range lineups {
		waitGroup.Add(1)

		go func(stats sportmonks.LineupPlayer) {
			p.PlayerProcessor.ProcessPlayerStats(&stats, bench)
			defer waitGroup.Done()
		}(player)
	}
}

func (p Processor) handleGoalEvents(g []sportmonks.GoalEvent) {
	for _, goal := range g {
		waitGroup.Add(1)

		go func(e sportmonks.GoalEvent) {
			p.EventProcessor.ProcessGoalEvent(&e)
			defer waitGroup.Done()
		}(goal)
	}
}

func (p Processor) handleSubstitutionEvents(s []sportmonks.SubstitutionEvent) {
	for _, sub := range s {
		waitGroup.Add(1)

		go func(e sportmonks.SubstitutionEvent) {
			p.EventProcessor.ProcessSubstitutionEvent(&e)
			defer waitGroup.Done()
		}(sub)
	}
}
