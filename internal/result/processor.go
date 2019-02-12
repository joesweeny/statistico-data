package result

import (
	"github.com/joesweeny/statshub/internal/fixture"
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"github.com/joesweeny/statshub/internal/event"
	"github.com/joesweeny/statshub/internal/stats"
	"sync"
)

const result = "result"
const callLimit = 1800
var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	Repository
	FixtureRepo fixture.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
	PlayerProcessor stats.PlayerProcessor
	TeamProcessor stats.TeamProcessor
	EventProcessor event.Processor
}

func (p Processor) Process(command string, done chan bool) {
	switch command {
	case result:
		go p.allResults(done)
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

	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseResults(results, done)
}

func (p Processor) callClient(ids []int, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"lineup,bench,stats,goals,substitutions"}

	for _, id := range ids {
		if _, err := p.Repository.GetByFixtureId(id); err != ErrNotFound {
			continue
		}

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
	go p.handleTeams(fix.TeamStats.Data)
	go p.handlePlayers(fix.Lineup.Data, false)
	go p.handlePlayers(fix.Bench.Data, true)
	go p.handleGoalEvents(fix.Goals.Data)
	go p.handleSubstitutionEvents(fix.Subs.Data)

	created := p.Factory.createResult(&fix)

	if err := p.Repository.Insert(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Result struct: %+v\n,", err.Error(), created)
	}
}

func (p Processor) handleTeams(t []sportmonks.TeamStats) {
	for _, team := range t {
		waitGroup.Add(1)

		go func(s sportmonks.TeamStats) {
			p.TeamProcessor.ProcessTeamStats(&s)
			defer waitGroup.Done()
		}(team)
	}
}

func (p Processor) handlePlayers(lineups []sportmonks.LineupPlayer, bench bool) {
	for _, player := range lineups {
		waitGroup.Add(1)

		go func(s sportmonks.LineupPlayer) {
			p.PlayerProcessor.ProcessPlayerStats(&s, bench)
			defer waitGroup.Done()
		}(player)
	}
}

func (p Processor) handleGoalEvents(g []sportmonks.GoalEvent) {
	for _, goal := range g {
		waitGroup.Add(1)

		go func(s sportmonks.GoalEvent) {
			p.EventProcessor.ProcessGoalEvent(&s)
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
