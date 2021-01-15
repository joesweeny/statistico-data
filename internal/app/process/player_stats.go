package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
)

const playerStats = "player-stats"
const playerStatsCurrentSeason = "player-stats:current-season"
const playerStatsBySeasonId = "player-stats:by-season-id"

type PlayerStatsProcessor struct {
	playerStatsRepo app.PlayerStatsRepository
	seasonRepo     app.SeasonRepository
	requester       app.PlayerStatRequester
	clock           clockwork.Clock
	logger          *logrus.Logger
}

func (p PlayerStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case playerStats:
		go p.processAllSeasons(done)
	case playerStatsCurrentSeason:
		go p.processCurrentSeason(done)
	case playerStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.processSeason(uint64(id), done)
	default:
		p.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p PlayerStatsProcessor) processAllSeasons(done chan bool) {
	ids, err := p.seasonRepo.IDs()

	if err != nil {
		p.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := p.requester.PlayerStatsBySeasonIDs(ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processCurrentSeason(done chan bool) {
	ids, err := p.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		p.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := p.requester.PlayerStatsBySeasonIDs(ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := p.requester.PlayerStatsBySeasonIDs([]uint64{seasonID})

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
}

func NewPlayerStatsProcessor(
	r app.PlayerStatsRepository,
	s app.SeasonRepository,
	q app.PlayerStatRequester,
	c clockwork.Clock,
	log *logrus.Logger,
) *PlayerStatsProcessor {
	return &PlayerStatsProcessor{playerStatsRepo: r, seasonRepo: s, requester: q, clock: c, logger: log}
}
