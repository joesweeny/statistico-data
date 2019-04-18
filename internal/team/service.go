package team

import (
	"github.com/statistico/statistico-data/internal/fixture"
	"log"
	pb "github.com/statistico/statistico-data/internal/proto/team"
)

type Service struct {
	FixtureRepo fixture.Repository
	Logger 		*log.Logger
}

func (s Service) GetTeamsForSeason(r *pb.SeasonRequest, stream pb.TeamService_GetTeamsForSeasonServer) error {
	// Get Ids for teams for a given season

	// For each Id pass to handler to fetch Team and hydrate new Team proto struct then stream back

	// If fetching Team fails then log fatal error and return error to consumer
	return nil
}
