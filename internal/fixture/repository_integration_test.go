package fixture

import (
	"database/sql"
	"testing"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresFixtureRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newFixture(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_fixture")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})

	t.Run("returns error when ID primary key violates unique constraint", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newFixture(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})

	conn.Close()
}

func TestGetById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresFixtureRepository{Connection: conn}

	t.Run("fixture can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixture(43)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(43, f.ID)
		a.Equal(14567, f.SeasonID)
		a.Equal(165789, *f.RoundID)
		a.Nil(f.VenueID)
		a.Equal(451, f.HomeTeamID)
		a.Equal(924, f.AwayTeamID)
		a.Nil(f.RefereeID)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.Date.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
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
		_, err := db.Exec("delete from sportmonks_fixture")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newFixture(id int) *model.Fixture {
	var roundId = 165789

	return &model.Fixture{
		ID:         id,
		SeasonID:   14567,
		RoundID:    &roundId,
		HomeTeamID: 451,
		AwayTeamID: 924,
		Date:       time.Unix(1548086929, 0),
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}
}