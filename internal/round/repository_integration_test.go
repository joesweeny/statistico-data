package round

import (
	"github.com/joesweeny/statshub/internal/config"
	"database/sql"
	"testing"
	"fmt"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresRoundRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newRound(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_round")

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
		c := newRound(50)

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
	repo := PostgresRoundRepository{Connection: conn}

	t.Run("round can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newRound(43)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(43, r.ID)
		a.Equal("5", r.Name)
		a.Equal(4387, r.SeasonID)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.StartDate.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.EndDate.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if round does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.GetById(99)

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

func getConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		_, err := db.Exec("delete from sportmonks_round")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newRound(id int) *model.Round {
	return &model.Round{
		ID:        id,
		Name:      "5",
		SeasonID:  4387,
		StartDate: time.Unix(1548086929, 0),
		EndDate:   time.Unix(1548086929, 0),
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}