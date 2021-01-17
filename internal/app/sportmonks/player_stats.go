package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/helpers"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"sync"
	"time"
)

type PlayerStatsRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (p PlayerStatsRequester) PlayerStatsByFixtureIDs(ids []uint64) <-chan *app.PlayerStats {
	ch := make(chan *app.PlayerStats, 100)

	go p.parseStats(ids, ch)

	return ch
}

func (p PlayerStatsRequester) PlayerStatsBySeasonIDs(seasonIDs []uint64) <-chan *app.PlayerStats {
	ch := make(chan *app.PlayerStats, 100)

	go p.parseBySeasonIDs(seasonIDs, ch)

	return ch
}

func (p PlayerStatsRequester) PlayerStatsByDate(date time.Time, competitionIDS []uint64) <-chan *app.PlayerStats {
	ch := make(chan *app.PlayerStats, 1000)

	go p.parseByDate(competitionIDS, date, ch)

	return ch
}

func (p PlayerStatsRequester) parseBySeasonIDs(seasonIDs []uint64, ch chan<- *app.PlayerStats) {
	defer close(ch)

	wg := sync.WaitGroup{}

	for _, id := range seasonIDs {
		wg.Add(1)
		go p.sendSeasonRequest(id, ch, &wg)
	}

	wg.Wait()
}

func (p PlayerStatsRequester) parseByDate(competitionIDS []uint64, date time.Time, ch chan<- *app.PlayerStats) {
	defer close(ch)

	wg := sync.WaitGroup{}

	for _, competitionID := range competitionIDS {
		wg.Add(1)
		go p.sendByDateRequest(competitionID, date, ch, &wg)
	}

	wg.Wait()
}

func (p PlayerStatsRequester) parseStats(ids []uint64, ch chan<- *app.PlayerStats) {
	defer close(ch)

	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go p.sendStatsRequest(id, ch, &wg)
	}

	wg.Wait()
}

func (p PlayerStatsRequester) sendByDateRequest(seasonID uint64, date time.Time, ch chan<- *app.PlayerStats, wg *sync.WaitGroup) {
	results, _, err := p.client.FixturesByDate(
		context.Background(),
		date,
		[]string{"lineup", "bench"},
		map[string][]int{"leagues": {int(seasonID)}},
	)

	if err != nil {
		p.logger.Errorf(
			"Error when calling client '%s' when making season fixtures request. Season ID %d",
			err.Error(),
			seasonID,
		)

		wg.Done()
		return
	}

	for _, res := range results {
		for _, stats := range res.Lineups() {
			ch <- transformPlayerStats(&stats, false)
		}

		for _, stats := range res.Bench() {
			ch <- transformPlayerStats(&stats, true)
		}
	}

	wg.Done()
}

func (p PlayerStatsRequester) sendSeasonRequest(seasonID uint64, ch chan<- *app.PlayerStats, wg *sync.WaitGroup) {
	res, _, err := p.client.SeasonByID(context.Background(), int(seasonID), []string{"results.lineup", "results.bench"})

	if err != nil {
		p.logger.Errorf(
			"Error when calling client '%s' when making season fixtures request. Season ID %d",
			err.Error(),
			seasonID,
		)

		wg.Done()
		return
	}

	for _, res := range res.Results() {
		for _, stats := range res.Lineups() {
			ch <- transformPlayerStats(&stats, false)
		}

		for _, stats := range res.Bench() {
			ch <- transformPlayerStats(&stats, true)
		}
	}

	wg.Done()
}

func (p PlayerStatsRequester) sendStatsRequest(id uint64, ch chan<- *app.PlayerStats, wg *sync.WaitGroup) {
	var filters map[string][]int

	res, _, err := p.client.FixtureByID(context.Background(), int(id), []string{"lineup", "bench"}, filters)

	if err != nil {
		p.logger.Errorf(
			"Error when calling client '%s' when making fixtures request to parse player stats. Fixture ID %d",
			err.Error(),
			id,
		)

		wg.Done()
		return
	}

	for _, stats := range res.Lineups() {
		ch <- transformPlayerStats(&stats, false)
	}

	for _, stats := range res.Bench() {
		ch <- transformPlayerStats(&stats, true)
	}

	wg.Done()
}

func transformPlayerStats(s *spClient.PlayerStats, sub bool) *app.PlayerStats {
	return &app.PlayerStats{
		FixtureID:         uint64(s.FixtureID),
		PlayerID:          uint64(s.PlayerID),
		TeamID:            uint64(s.TeamID),
		Position:          s.Position,
		FormationPosition: s.FormationPosition,
		IsSubstitute:      sub,
		PlayerShots:       handlePlayerShots(&s.Stats.Shots),
		PlayerGoals:       handlePlayerGoals(&s.Stats.Goals),
		PlayerFouls:       handlePlayerFouls(&s.Stats.Fouls),
		YellowCards:       s.Stats.Cards.YellowCards,
		RedCard:           s.Stats.Cards.RedCards,
		PlayerCrosses:     handlePlayerCrosses(&s.Stats.Passing),
		PlayerPasses:      handlePlayerPasses(&s.Stats.Passing),
		Assists:           s.Stats.Goals.Assist,
		Offsides:          s.Stats.Other.Offsides,
		Saves:             s.Stats.Other.Saves,
		PlayerPenalties:   handlePenalties(&s.Stats.Other),
		HitWoodwork:       s.Stats.Other.HitWoodwork,
		Tackles:           s.Stats.Other.Tackles,
		Blocks:            s.Stats.Other.Blocks,
		Interceptions:     s.Stats.Other.Interceptions,
		Clearances:        s.Stats.Other.Clearances,
		MinutesPlayed:     s.Stats.Other.MinutesPlayed,
	}
}

func handlePlayerShots(s *spClient.Shots) app.PlayerShots {
	return app.PlayerShots{
		Total:  s.Total,
		OnGoal: s.OnGoal,
	}
}

func handlePlayerGoals(s *spClient.Goals) app.PlayerGoals {
	return app.PlayerGoals{
		Scored:   s.Scored,
		Conceded: s.Conceded,
	}
}

func handlePlayerFouls(s *spClient.Fouls) app.PlayerFouls {
	return app.PlayerFouls{
		Drawn:     s.Drawn,
		Committed: helpers.ParseNullableInt(s.Committed),
	}
}

func handlePlayerCrosses(s *spClient.MatchPasses) app.PlayerCrosses {
	return app.PlayerCrosses{
		Total:    s.TotalCrosses,
		Accuracy: s.CrossesAccuracy,
	}
}

func handlePlayerPasses(s *spClient.MatchPasses) app.PlayerPasses {
	return app.PlayerPasses{
		Total:    s.Passes,
		Accuracy: s.PassesAccuracy,
	}
}

func handlePenalties(s *spClient.AdditionalPlayerMatchStats) app.PlayerPenalties {
	return app.PlayerPenalties{
		Scored:    s.PenScored,
		Missed:    s.PenMissed,
		Saved:     s.PenSaved,
		Committed: s.PenCommitted,
		Won:       s.PenWon,
	}
}

func NewPlayerStatsRequester(client *spClient.HTTPClient, log *logrus.Logger) *PlayerStatsRequester {
	return &PlayerStatsRequester{client: client, logger: log}
}
