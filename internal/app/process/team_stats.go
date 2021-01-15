package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
)

const teamStats = "team-stats"
const teamStatsCurrentSeason = "team-stats:current-season"
const teamStatsBySeasonId = "team-stats:by-season-id"

type TeamStatsProcessor struct {
	teamStatsRepo app.TeamStatsRepository
	seasonRepo    app.SeasonRepository
	requester     app.TeamStatsRequester
	clock         clockwork.Clock
	logger        *logrus.Logger
}

func (t TeamStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case teamStats:
		go t.processAllSeasons(done)
	case teamStatsCurrentSeason:
		go t.processCurrentSeason(done)
	case teamStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go t.processSeason(uint64(id), done)
	default:
		t.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (t TeamStatsProcessor) processAllSeasons(done chan bool) {
	ids, err := t.seasonRepo.IDs()

	if err != nil {
		t.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := t.requester.TeamStatsBySeasonIDs(ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processCurrentSeason(done chan bool) {
	ids, err := t.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		t.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := t.requester.TeamStatsBySeasonIDs(ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := t.requester.TeamStatsBySeasonIDs([]uint64{seasonID})

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) persistStats(ch <-chan *app.TeamStats, done chan bool) {
	for stats := range ch {
		t.persist(stats)
	}

	done <- true
}

func (t TeamStatsProcessor) persist(x *app.TeamStats) {
	_, err := t.teamStatsRepo.ByFixtureAndTeam(x.FixtureID, x.TeamID)

	if err != nil {
		if err := t.teamStatsRepo.InsertTeamStats(x); err != nil {
			t.logger.Warningf("Error '%s' occurred when inserting team stats struct: %+v\n,", err.Error(), *x)
		}

		return
	}
}

func NewTeamStatsProcessor(
	r app.TeamStatsRepository,
	s app.SeasonRepository,
	q app.TeamStatsRequester,
	c clockwork.Clock,
	log *logrus.Logger,
) *TeamStatsProcessor {
	return &TeamStatsProcessor{
		teamStatsRepo: r,
		seasonRepo: s,
		requester: q,
		clock: c,
		logger: log,
	}
}
