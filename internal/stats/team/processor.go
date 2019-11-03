package team_stats

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const teamStatsByResultId = "team-stats:by-result-id"
const teamStatsBySeasonId = "team-stats:by-season-id"
const teamStatsToday = "team-stats:today"
const callLimit = 1500

var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	TeamRepository
	TeamFactory
	FixtureRepo app.FixtureRepository
	Logger      *log.Logger
	Client      *sportmonks.Client
}

func (p Processor) Process(command string, option string, done chan bool) {
	switch command {
	case teamStatsByResultId:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go p.byId(done, id)
		}
	case teamStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.bySeasonId(done, id)
	case teamStatsToday:
		go p.statsToday(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) byId(done chan bool, id int) {
	fix, err := p.FixtureRepo.ByID(uint64(id))

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []uint64{fix.ID}

	go p.processStats(ids, done)
}

func (p Processor) bySeasonId(done chan bool, id int) {
	// Adding a Clock.Now() here is a bit hacky. Redo by dynamically handling this
	fix, err := p.FixtureRepo.BySeasonID(uint64(id), p.Clock.Now())

	if err != nil {
		p.Logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", id, err.Error())
		return
	}

	var ids []uint64

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	go p.processStats(ids, done)
}

func (p Processor) statsToday(done chan bool) {
	now := p.Clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := p.FixtureRepo.IDsBetween(from, to)

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processStats(ids, done)
}

func (p Processor) processStats(ids []uint64, done chan bool) {
	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseStats(results, done)
}

func (p Processor) callClient(ids []uint64, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"stats"}

	for _, id := range ids {
		if *c >= callLimit {
			p.Logger.Printf("Api call limited reached %d calls\n", *c)
			break
		}

		res, err := p.Client.FixtureById(int(id), q, 5)

		*c++

		if err != nil {
			p.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		ch <- res.Data
	}

	close(ch)
}

func (p Processor) parseStats(ch <-chan sportmonks.Fixture, done chan bool) {
	for x := range ch {
		p.handleStats(x)
	}

	waitGroup.Wait()

	done <- true
}

func (p Processor) handleStats(fix sportmonks.Fixture) {
	waitGroup.Add(1)

	for _, x := range fix.TeamStats.Data {
		p.processTeamStats(&x)
	}

	defer waitGroup.Done()
}

func (p Processor) processTeamStats(s *sportmonks.TeamStats) {
	x, err := p.TeamRepository.ByFixtureAndTeam(uint64(s.FixtureID), uint64(s.TeamID))

	if err == ErrNotFound {
		created := p.TeamFactory.createTeamStats(s)

		if err := p.TeamRepository.InsertTeamStats(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Team Stats struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := p.TeamFactory.updateTeamStats(s, x)

	if err := p.TeamRepository.UpdateTeamStats(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Team Stats struct: %+v\n,", err.Error(), updated)
	}

	return
}
