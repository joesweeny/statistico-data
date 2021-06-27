package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-football-data/internal/app"
	"strconv"
	"time"
)

const teamStatsByDate = "team-stats:by-date"
const teamStatsBySeasonId = "team-stats:by-season-id"
const teamStatsByCompetitionId = "team-stats:by-competition-id"

type TeamStatsProcessor struct {
	teamStatsRepo app.TeamStatsRepository
	competitionRepo    app.CompetitionRepository
	seasonRepo    app.SeasonRepository
	requester     app.TeamStatsRequester
	clock         clockwork.Clock
	logger        *logrus.Logger
}

func (t TeamStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case teamStatsByDate:
		go t.processByDate(option, done)
	case teamStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go t.processSeason(uint64(id), done)
	case teamStatsByCompetitionId:
		id, _ := strconv.Atoi(option)
		go t.processCompetition(uint64(id), done)
	default:
		t.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (t TeamStatsProcessor) processByDate(date string, done chan bool) {
	d, err := time.Parse("2006-01-02", date)

	if err != nil {
		t.logger.Fatalf("Error parsing date in team stats processor: %s", err.Error())
		return
	}

	ids, err := t.competitionRepo.IDs()

	if err != nil {
		t.logger.Fatalf("Error fetching competition IDs in team stats processor: %s", err.Error())
		return
	}

	ch := t.requester.TeamStatsByDate(d, ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := t.requester.TeamStatsBySeasonIDs([]uint64{seasonID})

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processCompetition(competitionID uint64, done chan bool) {
	seasons, err := t.seasonRepo.ByCompetitionId(competitionID, "name_asc")

	if err != nil {
		t.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	var ids []uint64

	for _, season := range seasons {
		ids = append(ids, season.ID)
	}

	ch := t.requester.TeamStatsBySeasonIDs(ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) persistStats(ch <-chan app.TeamStats, done chan bool) {
	for stats := range ch {
		t.persist(stats)
	}

	done <- true
}

func (t TeamStatsProcessor) persist(x app.TeamStats) {
	_, err := t.teamStatsRepo.ByFixtureAndTeam(x.FixtureID, x.TeamID)

	if err != nil {
		if err := t.teamStatsRepo.InsertTeamStats(&x); err != nil {
			t.logger.Errorf("Error '%s' occurred when inserting team stats struct: %+v\n,", err.Error(), x)
		}

		return
	}

	if err := t.teamStatsRepo.UpdateTeamStats(&x); err != nil {
		t.logger.Errorf("Error '%s' occurred when updating team stats struct: %+v\n,", err.Error(), x)
	}

	return
}

func NewTeamStatsProcessor(
	r app.TeamStatsRepository,
	c app.CompetitionRepository,
	s app.SeasonRepository,
	q app.TeamStatsRequester,
	cl clockwork.Clock,
	log *logrus.Logger,
) *TeamStatsProcessor {
	return &TeamStatsProcessor{
		teamStatsRepo: r,
		competitionRepo: c,
		seasonRepo: s,
		requester: q,
		clock: cl,
		logger: log,
	}
}
