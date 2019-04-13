package team_stats

import (
	"database/sql"
	"fmt"
	"github.com/statistico/statistico-data/internal/config"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTeamPlayerStats(t *testing.T) {
	conn, cleanUp := getTeamConnection(t)
	repo := PostgresTeamStatsRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newTeamStats(42, 65)

			if err := repo.InsertTeamStats(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_team_stats")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})
}

func TestByFixtureAndTeam(t *testing.T) {
	conn, cleanUp := getTeamConnection(t)
	repo := PostgresTeamStatsRepository{Connection: conn}

	t.Run("team stats can be retrieved by fixture and team IDs", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newTeamStats(42, 65)

		if err := repo.InsertTeamStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByFixtureAndTeam(42, 65)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(42, m.FixtureID)
		a.Equal(65, m.TeamID)
		a.Nil(m.TeamShots.Total)
		a.Nil(m.TeamShots.OnGoal)
		a.Nil(m.TeamShots.OffGoal)
		a.Nil(m.TeamShots.Blocked)
		a.Nil(m.TeamShots.InsideBox)
		a.Nil(m.TeamShots.OutsideBox)
		a.Nil(m.TeamPasses.Total)
		a.Nil(m.TeamPasses.Accuracy)
		a.Nil(m.TeamPasses.Percentage)
		a.Nil(m.TeamAttacks.Total)
		a.Nil(m.TeamAttacks.Dangerous)
		a.Nil(m.Fouls)
		a.Nil(m.Corners)
		a.Nil(m.Offsides)
		a.Nil(m.Possession)
		a.Nil(m.YellowCards)
		a.Nil(m.RedCards)
		a.Nil(m.Saves)
		a.Nil(m.Substitutions)
		a.Nil(m.GoalKicks)
		a.Nil(m.GoalAttempts)
		a.Nil(m.FreeKicks)
		a.Nil(m.ThrowIns)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if stats does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.ByFixtureAndTeam(99, 100)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %s", ErrNotFound, err)
		}
	})

	conn.Close()
}

func TestUpdateTeamStats(t *testing.T) {
	conn, cleanUp := getTeamConnection(t)
	repo := PostgresTeamStatsRepository{Connection: conn}

	t.Run("modifies existing team stats record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newTeamStats(42, 65)

		if err := repo.InsertTeamStats(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var shotTotal = 10
		var shotOnGoal = 2
		var shotOffGoal = 1
		var shotBlocked = 5
		var shotInside = 8
		var shotOutside = 2
		var passTotal = 156
		var passAcc = 78
		var passPer = 98
		var attTotal = 50
		var attDan = 50
		var fouls = 56
		var corner = 4
		var offside = 3
		var poss = 56
		var yellow = 4
		var red = 0
		var save = 0
		var goalKicks = 2
		var goalAttempt = 2
		var throwsIns = 9
		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		m.TeamShots.Total = &shotTotal
		m.TeamShots.OnGoal = &shotOnGoal
		m.TeamShots.OffGoal = &shotOffGoal
		m.TeamShots.Blocked = &shotBlocked
		m.TeamShots.InsideBox = &shotInside
		m.TeamShots.OutsideBox = &shotOutside
		m.TeamPasses.Total = &passTotal
		m.TeamPasses.Accuracy = &passAcc
		m.TeamPasses.Percentage = &passPer
		m.TeamAttacks.Total = &attTotal
		m.TeamAttacks.Dangerous = &attDan
		m.Fouls = &fouls
		m.Corners = &corner
		m.Offsides = &offside
		m.Possession = &poss
		m.YellowCards = &yellow
		m.RedCards = &red
		m.Saves = &save
		m.GoalKicks = &goalKicks
		m.GoalAttempts = &goalAttempt
		m.ThrowIns = &throwsIns
		m.UpdatedAt = d

		if err := repo.UpdateTeamStats(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByFixtureAndTeam(42, 65)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(42, m.FixtureID)
		a.Equal(65, m.TeamID)
		a.Equal(10, *m.TeamShots.Total)
		a.Equal(2, *m.TeamShots.OnGoal)
		a.Equal(1, *m.TeamShots.OffGoal)
		a.Equal(5, *m.TeamShots.Blocked)
		a.Equal(8, *m.TeamShots.InsideBox)
		a.Equal(2, *m.TeamShots.OutsideBox)
		a.Equal(156, *m.TeamPasses.Total)
		a.Equal(78, *m.TeamPasses.Accuracy)
		a.Equal(98, *m.TeamPasses.Percentage)
		a.Equal(50, *m.TeamAttacks.Total)
		a.Equal(50, *m.TeamAttacks.Dangerous)
		a.Equal(56, *m.Fouls)
		a.Equal(4, *m.Corners)
		a.Equal(3, *m.Offsides)
		a.Equal(56, *m.Possession)
		a.Equal(4, *m.YellowCards)
		a.Equal(0, *m.RedCards)
		a.Equal(0, *m.Saves)
		a.Nil(m.Substitutions)
		a.Equal(2, *m.GoalKicks)
		a.Equal(2, *m.GoalAttempts)
		a.Nil(m.FreeKicks)
		a.Equal(9, *m.ThrowIns)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if stats does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		err := repo.UpdateTeamStats(newTeamStats(1, 2))

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %v", ErrNotFound, err)
		}
	})

	conn.Close()
}

var db = config.GetConfig().Database

var teamDb = config.GetConfig().Database

func getTeamConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		teamDb.Host, teamDb.Port, teamDb.User, teamDb.Password, teamDb.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		_, err := db.Exec("delete from sportmonks_team_stats")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newTeamStats(fixtureId, teamId int) *model.TeamStats {
	return &model.TeamStats{
		FixtureID:   fixtureId,
		TeamID:      teamId,
		TeamShots:   model.TeamShots{},
		TeamPasses:  model.TeamPasses{},
		TeamAttacks: model.TeamAttacks{},
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}
