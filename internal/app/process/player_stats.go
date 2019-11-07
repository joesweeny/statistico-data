package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"strings"
	"time"
)

const playerStatsByResultId = "player-stats:by-result-id"
const playerStatsBySeasonId = "player-stats:by-season-id"
const playerStatsToday = "player-stats:today"

type PlayerStatsProcessor struct {
	playerStatsRepo app.PlayerStatsRepository
	fixtureRepo app.FixtureRepository
	requester app.PlayerStatRequester
	clock clockwork.Clock
	logger     *logrus.Logger
}

func (p PlayerStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case playerStatsByResultId:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go p.processByID(done, uint64(id))
		}
	case playerStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.processSeason(done, uint64(id))
	case playerStatsToday:
		go p.processToday(done)
	default:
		p.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p PlayerStatsProcessor) processByID(done chan bool, id uint64) {
	fix, err := p.fixtureRepo.ByID(id)

	if err != nil {
		p.logger.Fatalf("Error when retrieving fixtures for ID: %d, %s", id, err.Error())
		return
	}

	ch := p.requester.PlayerStatsByFixtureIDs([]uint64{fix.ID})

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processSeason(done chan bool, seasonID uint64) {
	fix, err := p.fixtureRepo.BySeasonID(seasonID)

	if err != nil {
		p.logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", seasonID, err.Error())
		return
	}

	var ids []uint64

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	ch := p.requester.PlayerStatsByFixtureIDs(ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processToday(done chan bool) {
	now := p.clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := p.fixtureRepo.IDsBetween(from, to)

	if err != nil {
		p.logger.Fatalf("Error when retrieving fixture ids in player stats processor: %s", err.Error())
		return
	}

	ch := p.requester.PlayerStatsByFixtureIDs(ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) persistStats(ch <-chan *app.PlayerStats, done chan bool) {
	for stats := range ch {
		p.persist(stats)
	}

	done <- true
}

func (p PlayerStatsProcessor) persist(x *app.PlayerStats) {
	_, err := p.playerStatsRepo.ByFixtureAndPlayer(x.FixtureID, x.PlayerID)

	if err != nil {
		if err := p.playerStatsRepo.Insert(x); err != nil {
			p.logger.Warningf("Error '%s' occurred when inserting player stats struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := p.playerStatsRepo.Update(x); err != nil {
		p.logger.Warningf("Error '%s' occurred when updating player stats struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewPlayerStatsProcessor(
	r app.PlayerStatsRepository,
	f app.FixtureRepository,
	q app.PlayerStatRequester,
	c clockwork.Clock,
	log *logrus.Logger,
) *PlayerStatsProcessor {
	return &PlayerStatsProcessor{playerStatsRepo: r, fixtureRepo: f, requester: q, clock: c, logger: log}
}
