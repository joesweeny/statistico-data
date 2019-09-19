package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	now   = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
	clock = clockwork.NewFakeClockAt(now)
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := CountryRepository{connection: conn, clock: clock}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newCountry(i)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_country")

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
		c := newCountry(10)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestUpdate(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := CountryRepository{connection: conn, clock: clock}

	t.Run("modifies existing record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCountry(100)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		s := "Scotland"
		c.Name = s

		if err := repo.Update(c); err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.GetById(c.ID)

		if err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		got := r.Name
		want := s

		assert.Equal(t, got, want)
	})

	t.Run("returns error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCountry(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func TestGetById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := CountryRepository{connection: conn, clock: clock}

	t.Run("country can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		c := newCountry(62)

		if err := repo.Insert(c); err != nil {
			t.Fatalf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(62)

		if err != nil {
			t.Fatalf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(62, r.ID)
		a.Equal("England", r.Name)
		a.Equal("Europe", r.Continent)
		a.Equal("ENG", r.ISO)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if country does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.GetById(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func getConnection(t *testing.T) (*sql.DB, func()) {
	db := config.GetConfig().Database

	dsn := "host=%s port=%s user=%s " + "password=%s dbname=%s sslmode=disable"

	psqlInfo := fmt.Sprintf(dsn, db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return conn, func() {
		_, err := conn.Exec("delete from sportmonks_country")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newCountry(id int) *app.Country {
	c := app.Country{
		ID:        id,
		Name:      "England",
		Continent: "Europe",
		ISO:       "ENG",
	}

	return &c
}
