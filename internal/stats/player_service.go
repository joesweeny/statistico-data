package stats

import (
	"github.com/statistico/statistico-data/internal/fixture"
	pb "github.com/statistico/statistico-data/internal/proto/stats"
	"log"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	PlayerRepository PlayerRepository
	FixtureRepo      fixture.Repository
	Logger 			 *log.Logger
}

func (s Service) GetPlayerStatsForFixture(c context.Context, r *pb.FixtureRequest) (*pb.FixtureResponse, error) {
	// Get fixture by ID
	fix, err := s.FixtureRepo.ById(r.FixtureId)
	// Get Player stats slice for Home Team
	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}
	// Get Player stats slice for away team
	res := pb.FixtureResponse{}

	home := s.PlayerRepository.
	// Pass home slice through handler to hydrate proto stats structs and assign to response

	// Pass away slice through handler to hydrate proto stats structs and assign to response

	return &res, nil
}
