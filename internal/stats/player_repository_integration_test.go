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