package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-football-data/internal/app"
)

const squad = "squad"
const squadCurrentSeason = "squad:current-season"

type SquadProcessor struct {
	squadRepo  app.SquadRepository
	seasonRepo app.SeasonRepository
	requester  app.SquadRequester
	logger     *logrus.Logger
}

func (s SquadProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case squad:
		go s.processAllSeasons(done)
	case squadCurrentSeason:
		go s.processCurrentSeason(done)
	default:
		s.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (s SquadProcessor) processAllSeasons(done chan bool) {
	ids, err := s.seasonRepo.IDs()

	if err != nil {
		s.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := s.requester.SquadsBySeasonIDs(ids)

	go s.persistSquads(ch, done)
}

func (s SquadProcessor) processCurrentSeason(done chan bool) {
	ids, err := s.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		s.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := s.requester.SquadsBySeasonIDs(ids)

	go s.persistSquads(ch, done)
}

func (s SquadProcessor) persistSquads(ch <-chan *app.Squad, done chan bool) {
	for squad := range ch {
		s.persist(squad)
	}

	done <- true
}

func (s SquadProcessor) persist(x *app.Squad) {
	_, err := s.squadRepo.BySeasonAndTeam(x.TeamID, x.SeasonID)

	if err != nil {
		if newErr := s.squadRepo.Insert(x); newErr != nil {
			s.logger.Warningf("Error '%s' occurred when inserting squad struct: %+v\n,", newErr.Error(), *x)
		}

		return
	}

	if newErr := s.squadRepo.Update(x); newErr != nil {
		s.logger.Warningf("Error '%s' occurred when updating squad struct: %+v\n,", newErr.Error(), *x)
	}

	return
}

func NewSquadProcessor(s app.SquadRepository, a app.SeasonRepository, r app.SquadRequester, log *logrus.Logger) *SquadProcessor {
	return &SquadProcessor{squadRepo: s, seasonRepo: a, requester: r, logger: log}
}
