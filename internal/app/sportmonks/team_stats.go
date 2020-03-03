package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/helpers"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
)

type TeamStatsRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (t TeamStatsRequester) TeamStatsByFixtureIDs(ids []uint64) <-chan *app.TeamStats {
	ch := make(chan *app.TeamStats, 100)

	go t.parseStats(ids, ch)

	return ch
}

func (t TeamStatsRequester) parseStats(ids []uint64, ch chan<- *app.TeamStats) {
	defer close(ch)

	var filters map[string][]int

	for _, id := range ids {
		res, _, err := t.client.FixtureByID(context.Background(), int(id), []string{"stats"}, filters)

		if err != nil {
			t.logger.Warnf(
				"Error when calling client '%s' when making fixtures request to parse team stats for fixture %d",
				err.Error(),
				id,
			)
			continue
		}

		for _, stats := range res.TeamStats() {
			ch <- transformTeamStats(&stats)
		}
	}
}

func transformTeamStats(s *spClient.TeamStats) *app.TeamStats {
	return &app.TeamStats{
		FixtureID:     uint64(s.FixtureID),
		TeamID:        uint64(s.TeamID),
		TeamShots:     handleTeamShots(&s.Shots),
		TeamPasses:    handleTeamPasses(&s.Passes),
		TeamAttacks:   handleTeamAttacks(&s.Attacks),
		Fouls:         helpers.ParseNullableInt(s.Fouls),
		Corners:       helpers.ParseNullableInt(s.Corners),
		Offsides:      helpers.ParseNullableInt(s.Offsides),
		Possession:    helpers.ParseNullableInt(s.PossessionTime),
		YellowCards:   helpers.ParseNullableInt(s.YellowCards),
		RedCards:      helpers.ParseNullableInt(s.RedCards),
		Saves:         helpers.ParseNullableInt(s.Saves),
		Substitutions: helpers.ParseNullableInt(s.Substitutions),
		GoalKicks:     helpers.ParseNullableInt(s.GoalKick),
		GoalAttempts:  helpers.ParseNullableInt(s.GoalAttempts),
		FreeKicks:     helpers.ParseNullableInt(s.FreeKick),
		ThrowIns:      helpers.ParseNullableInt(s.ThrowIn),
	}
}

func handleTeamShots(s *spClient.TeamShots) app.TeamShots {
	return app.TeamShots{
		Total:      helpers.ParseFlexInt(s.Total),
		OnGoal:     helpers.ParseFlexInt(s.OnGoal),
		OffGoal:    helpers.ParseFlexInt(s.OffGoal),
		Blocked:    helpers.ParseFlexInt(s.Blocked),
		InsideBox:  helpers.ParseFlexInt(s.InsideBox),
		OutsideBox: helpers.ParseFlexInt(s.OutsideBox),
	}
}

func handleTeamPasses(s *spClient.TeamPasses) app.TeamPasses {
	return app.TeamPasses{
		Total:      helpers.ParseFlexInt(s.Total),
		Accuracy:   helpers.ParseFlexInt(s.Accurate),
		Percentage: s.Percentage,
	}
}

func handleTeamAttacks(s *spClient.TeamAttacks) app.TeamAttacks {
	return app.TeamAttacks{
		Total:     helpers.ParseFlexInt(s.Total),
		Dangerous: helpers.ParseFlexInt(s.Dangerous),
	}
}

func NewTeamStatsRequester(client *spClient.HTTPClient, log *logrus.Logger) *TeamStatsRequester {
	return &TeamStatsRequester{client: client, logger: log}
}
