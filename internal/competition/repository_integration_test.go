package competition

import (
	"testing"
	"database/sql"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresCompetitionRepository{Connection: conn}

	t.Run("increases table count", func (t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newCompetition(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_competition")

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
		c := newCompetition(50)

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
	repo := PostgresCompetitionRepository{Connection: conn}

	t.Run("competition can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		c := newCompetition(45)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(45)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(45, r.ID)
		a.Equal("Premier League", r.Name)
		a.Equal(462, r.CountryID)
		a.Equal(false, r.IsCup)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if country does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.GetById(780);

		if err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}

		assert.Equal(t, err, ErrNotFound)
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
		_, err := db.Exec("delete from sportmonks_competition")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newCompetition(id int) model.Competition {
	return model.Competition{
		ID:         id,
		Name:       "Premier League",
		CountryID:  462,
		IsCup:      false,
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}
}