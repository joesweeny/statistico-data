package stats

import (
	"github.com/joesweeny/statshub/internal/config"
	"database/sql"
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/joesweeny/statshub/internal/model"
	"time"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getPlayerConnection(t)
	repo := PostgresPlayerStatsRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newPlayerStats(42, 65)

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

		m := newPlayerStats(30, 672)

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
		a.Equal("M", r.Position)
		a.Nil(r.FormationPosition)
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

func newPlayerStats(fixtureId, playerId int) *model.PlayerStats {
	return &model.PlayerStats{
		FixtureID: fixtureId,
		PlayerID: playerId,
		TeamID: 100,
		Position: "M",
		IsSubstitute: false,
		PlayerShots: model.PlayerShots{},
		PlayerGoals: model.PlayerGoals{},
		PlayerFouls: model.PlayerFouls{},
		PlayerCrosses: model.PlayerCrosses{},
		PlayerPasses: model.PlayerPasses{},
		PlayerPenalties: model.PlayerPenalties{},
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}