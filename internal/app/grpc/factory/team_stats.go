package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
)

type TeamStatsFactory struct {
	repo app.TeamStatsRepository
	logger 			*logrus.Logger
}

func (t TeamStatsFactory) BuildTeamStats(f *app.Fixture, teamID uint64) (*proto.TeamStats, error) {
	team, err := t.repo.ByFixtureAndTeam(f.ID, teamID)

	if err != nil {
		return nil, t.returnLoggedError(f.ID, err)
	}

	return teamStatsToProto(team), nil
}

func (t TeamStatsFactory) returnLoggedError(id uint64, err error) error {
	t.logger.Warnf("error when hydrating proto player stats: fixture %d. error %s", id, err.Error())
	return err
}

func NewTeamStatsFactory(r app.TeamStatsRepository, log *logrus.Logger) *TeamStatsFactory {
	return &TeamStatsFactory{repo: r, logger: log}
}
