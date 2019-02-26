package manager

import (
	"database/sql"
	"fmt"
	"github.com/joesweeny/statistico-data/internal/config"
	"github.com/joesweeny/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresManagerRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newManager(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_manager")

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
		c := newManager(50)

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
	repo := PostgresManagerRepository{Connection: conn}

	t.Run("manager can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newManager(10)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.Id(10)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(10, r.ID)
		a.Nil(r.TeamID)
		a.Equal(167, r.CountryID)
		a.Equal("Manuel", r.FirstName)
		a.Equal("Pellegrini", r.LastName)
		a.Equal("Chilean", r.Nationality)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if manager does not exist", func(t *testing.T) {
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

func TestUpdate(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresManagerRepository{Connection: conn}

	t.Run("modifies existing manager", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newManager(56)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var teamId = 574
		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		m.TeamID = &teamId
		m.UpdatedAt = d

		if err := repo.Update(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.Id(56)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(56, r.ID)
		a.Equal(574, *r.TeamID)
		a.Equal(167, r.CountryID)
		a.Equal("Manuel", r.FirstName)
		a.Equal("Pellegrini", r.LastName)
		a.Equal("Chilean", r.Nationality)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if manager does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newManager(146)

		err := repo.Update(c)

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

func getConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		_, err := db.Exec("delete from sportmonks_manager")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newManager(id int) *model.Manager {
	return &model.Manager{
		ID:          id,
		CountryID:   167,
		FirstName:   "Manuel",
		LastName:    "Pellegrini",
		Nationality: "Chilean",
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}
