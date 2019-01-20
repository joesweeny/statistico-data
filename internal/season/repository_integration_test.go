package season

import (
	"github.com/joesweeny/statshub/internal/config"
	"testing"
	"database/sql"
	"fmt"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	"github.com/stretchr/testify/assert"
	_ "github.com/lib/pq"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresSeasonRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			s := newSeason(i)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_season")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})
}

var db = config.GetConfig().Database

func getConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		_, err := db.Exec("delete from sportmonks_season")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newSeason(id int) model.Season {
	return model.Season{
		ID:         id,
		Name:       "2018-2019",
		LeagueID:   560,
		IsCurrent:  true,
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}
}
