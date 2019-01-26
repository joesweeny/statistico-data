package player

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
	repo := PostgresPlayerRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newPlayer(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_player")

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
		c := newPlayer(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})

	conn.Close()
}

func TestId(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresPlayerRepository{Connection: conn}

	t.Run("player can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayer(43)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.Id(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(43, r.ID)
		a.Equal(154, m.CountryId)
		a.Equal("Manuel", m.FirstName)
		a.Equal("Lanzini", m.LastName)
		a.Equal("Buenos Aires", *m.BirthPlace)
		a.Equal("1984-03-12", *m.DateOfBirth)
		a.Equal(3, m.PositionID)
		a.Equal("path/to/image", *m.Image)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if player does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.Id(99)

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
		_, err := db.Exec("delete from sportmonks_player")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newPlayer(id int) *model.Player {
	var place = "Buenos Aires"
	var dob = "1984-03-12"
	var path = "path/to/image"
	return &model.Player{
		ID:          id,
		CountryId:   154,
		FirstName:   "Manuel",
		LastName:    "Lanzini",
		BirthPlace:  &place,
		DateOfBirth: &dob,
		PositionID:  3,
		Image:       &path,
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}
