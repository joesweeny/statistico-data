package player_stats

import (
	"database/sql"
	"fmt"
	"github.com/statistico/statistico-data/internal/config"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getPlayerConnection(t)
	repo := PostgresPlayerStatsRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newPlayerStats(42, 65, 100, i)

			if err := repo.InsertPlayerStats(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_player_stats")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})
}

func TestByFixtureAndPlayer(t *testing.T) {
	conn, cleanUp := getPlayerConnection(t)
	repo := PostgresPlayerStatsRepository{Connection: conn}

	t.Run("player stats can be retrieved by fixture and player IDs", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayerStats(30, 672, 100, 1)

		if err := repo.InsertPlayerStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByFixtureAndPlayer(30, 672)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(30, r.FixtureID)
		a.Equal(672, r.PlayerID)
		a.Equal(100, r.TeamID)
		a.Equal("M", *r.Position)
		a.Equal(1, *r.FormationPosition)
		a.False(r.IsSubstitute)
		a.Nil(m.PlayerShots.Total)
		a.Nil(m.PlayerShots.OnGoal)
		a.Nil(m.PlayerGoals.Scored)
		a.Nil(m.PlayerGoals.Conceded)
		a.Nil(m.PlayerFouls.Drawn)
		a.Nil(m.PlayerFouls.Committed)
		a.Nil(m.YellowCards)
		a.Nil(m.RedCard)
		a.Nil(m.PlayerCrosses.Total)
		a.Nil(m.PlayerCrosses.Accuracy)
		a.Nil(m.PlayerPasses.Total)
		a.Nil(m.PlayerPasses.Accuracy)
		a.Nil(m.Assists)
		a.Nil(m.Offsides)
		a.Nil(m.Saves)
		a.Nil(m.PlayerPenalties.Scored)
		a.Nil(m.PlayerPenalties.Missed)
		a.Nil(m.PlayerPenalties.Saved)
		a.Nil(m.PlayerPenalties.Committed)
		a.Nil(m.PlayerPenalties.Won)
		a.Nil(m.HitWoodwork)
		a.Nil(m.Tackles)
		a.Nil(m.Blocks)
		a.Nil(m.Interceptions)
		a.Nil(m.Clearances)
		a.Nil(m.MinutesPlayed)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if stats does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.ByFixtureAndPlayer(99, 100)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %s", ErrNotFound, err)
		}
	})

	conn.Close()
}

func TestUpdatePlayerStats(t *testing.T) {
	conn, cleanUp := getPlayerConnection(t)
	repo := PostgresPlayerStatsRepository{Connection: conn}

	t.Run("modifies existing player stats record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayerStats(30, 672, 100, 1)

		if err := repo.InsertPlayerStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var formPos = 3
		var shotTotal = 4
		var shotOnGoal = 1
		var goalsTotal = 1
		var goalsCon = 0
		var foulsDrawn = 10
		var foulsCommitted = 4
		var yellow = 1
		var red = 0
		var crossTotal = 45
		var crossAccuracy = 60
		var passTotal = 68
		var passAccuracy = 90
		var assist = 3
		var offside = 3
		var saves = 0
		var penScored = 0
		var penMissed = 0
		var penSaved = 0
		var penCommitted = 0
		var penWon = 0
		var woodWork = 4
		var tackles = 8
		var blocks = 2
		var interceptions = 3
		var clearance = 1
		var minPlayed = 90
		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		m.IsSubstitute = true
		m.FormationPosition = &formPos
		m.PlayerShots.Total = &shotTotal
		m.PlayerShots.OnGoal = &shotOnGoal
		m.PlayerGoals.Scored = &goalsTotal
		m.PlayerGoals.Conceded = &goalsCon
		m.PlayerFouls.Drawn = &foulsDrawn
		m.PlayerFouls.Committed = &foulsCommitted
		m.YellowCards = &yellow
		m.RedCard = &red
		m.PlayerCrosses.Total = &crossTotal
		m.PlayerCrosses.Accuracy = &crossAccuracy
		m.PlayerPasses.Total = &passTotal
		m.PlayerPasses.Accuracy = &passAccuracy
		m.Assists = &assist
		m.Offsides = &offside
		m.Saves = &saves
		m.PlayerPenalties.Scored = &penScored
		m.PlayerPenalties.Missed = &penMissed
		m.PlayerPenalties.Saved = &penSaved
		m.PlayerPenalties.Committed = &penCommitted
		m.PlayerPenalties.Won = &penWon
		m.HitWoodwork = &woodWork
		m.Tackles = &tackles
		m.Blocks = &blocks
		m.Interceptions = &interceptions
		m.Clearances = &clearance
		m.MinutesPlayed = &minPlayed
		m.UpdatedAt = d

		if err := repo.UpdatePlayerStats(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByFixtureAndPlayer(30, 672)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(30, r.FixtureID)
		a.Equal(672, r.PlayerID)
		a.Equal(100, r.TeamID)
		a.Equal("M", *r.Position)
		a.Equal(3, *r.FormationPosition)
		a.True(r.IsSubstitute)
		a.Equal(4, *m.PlayerShots.Total)
		a.Equal(1, *m.PlayerShots.OnGoal)
		a.Equal(1, *m.PlayerGoals.Scored)
		a.Equal(0, *m.PlayerGoals.Conceded)
		a.Equal(10, *m.PlayerFouls.Drawn)
		a.Equal(4, *m.PlayerFouls.Committed)
		a.Equal(1, *m.YellowCards)
		a.Equal(0, *m.RedCard)
		a.Equal(45, *m.PlayerCrosses.Total)
		a.Equal(60, *m.PlayerCrosses.Accuracy)
		a.Equal(68, *m.PlayerPasses.Total)
		a.Equal(90, *m.PlayerPasses.Accuracy)
		a.Equal(3, *m.Assists)
		a.Equal(3, *m.Offsides)
		a.Equal(0, *m.Saves)
		a.Equal(0, *m.PlayerPenalties.Scored)
		a.Equal(0, *m.PlayerPenalties.Missed)
		a.Equal(0, *m.PlayerPenalties.Saved)
		a.Equal(0, *m.PlayerPenalties.Committed)
		a.Equal(0, *m.PlayerPenalties.Won)
		a.Equal(4, *m.HitWoodwork)
		a.Equal(8, *m.Tackles)
		a.Equal(2, *m.Blocks)
		a.Equal(3, *m.Interceptions)
		a.Equal(1, *m.Clearances)
		a.Equal(90, *m.MinutesPlayed)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if stats does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		err := repo.UpdatePlayerStats(newPlayerStats(1, 2, 100, 1))

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %v", ErrNotFound, err)
		}
	})

	conn.Close()
}

func TestByFixtureAndTeam(t *testing.T) {
	conn, cleanUp := getPlayerConnection(t)
	repo := PostgresPlayerStatsRepository{Connection: conn}

	t.Run("returns a slice of player stats structs", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayerStats(30, 6, 22, 4)

		if err := repo.InsertPlayerStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		for i := 1; i < 5; i++ {
			m := newPlayerStats(30, i, 100, i)

			if err := repo.InsertPlayerStats(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		stats, err := repo.ByFixtureAndTeam(30, 100)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t,4, len(stats))

		for _, s := range stats {
			assert.Equal(t, 100, s.TeamID)
		}
	})

	t.Run("player stats structs are ordered ascending by formation position", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayerStats(30, 16, 22, 10)

		if err := repo.InsertPlayerStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		for i := 5; i >= 1; i-- {
			m := newPlayerStats(30, i, 100, i)

			if err := repo.InsertPlayerStats(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		stats, err := repo.ByFixtureAndTeam(30, 100)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t,5, len(stats))

		for i, s := range stats {
			assert.Equal(t, 30, s.FixtureID)
			assert.Equal(t, 100, s.TeamID)
			assert.Equal(t, i + 1, *s.FormationPosition)
			assert.Equal(t, i + 1, s.PlayerID)
		}
	})
}

var db = config.GetConfig().Database

func getPlayerConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		_, err := db.Exec("delete from sportmonks_player_stats")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newPlayerStats(fixtureId, playerId, teamId, formation int) *model.PlayerStats {
	pos := "M"
	return &model.PlayerStats{
		FixtureID:       fixtureId,
		PlayerID:        playerId,
		TeamID:          teamId,
		Position:        &pos,
		FormationPosition: &formation,
		IsSubstitute:    false,
		PlayerShots:     model.PlayerShots{},
		PlayerGoals:     model.PlayerGoals{},
		PlayerFouls:     model.PlayerFouls{},
		PlayerCrosses:   model.PlayerCrosses{},
		PlayerPasses:    model.PlayerPasses{},
		PlayerPenalties: model.PlayerPenalties{},
		CreatedAt:       time.Unix(1546965200, 0),
		UpdatedAt:       time.Unix(1546965200, 0),
	}
}
