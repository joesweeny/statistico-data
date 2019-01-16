package country

import (
	"testing"
	"fmt"
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/satori/go.uuid"
	"time"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := postgresCountryRepository{conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		for i := 1; i < 4; i++ {
			var id = uuid.Must(uuid.NewV4(), nil)
			c := newCountry(id)

			if err := repo.insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from country")

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
		c := newCountry(uuid.UUID{})

		if err := repo.insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})

	conn.Close()
}

func TestUpdate(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := postgresCountryRepository{conn}

	t.Run("modifies existing record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		var id = uuid.Must(uuid.NewV4(), nil)
		c := newCountry(id)

		if err := repo.insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		s := "Scotland"
		c.Name = s

		if err := repo.update(c); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.getById(c.ID)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		got := r.Name
		want := s

		assert.Equal(t, got, want)
	})

	t.Run("returns error if record does not exist", func (t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCountry(uuid.NewV4())

		err := repo.update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})

	conn.Close()
}

func TestGetById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := postgresCountryRepository{conn}

	t.Run("country can be retrieved by ID", func (t *testing.T) {
		t.Helper()
		defer cleanUp()

		var id = uuid.Must(uuid.NewV4(), nil)
		c := newCountry(id)

		if err := repo.insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.getById(id)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal("England", r.Name)
		a.Equal("Europe", r.Continent)
		a.Equal("ENGLAND", r.ISO)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if country does not exist", func (t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.getById(uuid.UUID{}); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})

	conn.Close()
}

func TestGetByExternalId(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := postgresCountryRepository{conn}

	t.Run("country can be retrieved by External ID", func (t *testing.T) {
		t.Helper()
		defer cleanUp()

		var id = uuid.Must(uuid.NewV4(), nil)
		c := newCountry(id)

		if err := repo.insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.getByExternalId(c.ExternalID)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal("England", r.Name)
		a.Equal("Europe", r.Continent)
		a.Equal("ENG", r.ISO)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if country does not exist", func (t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.getByExternalId(999); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
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
		_, err := db.Exec("delete from country")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newCountry(u uuid.UUID) model.Country {
	c := model.Country{
		ID:         u,
		ExternalID: 1,
		Name:       "England",
		Continent:  "Europe",
		ISO:        "ENG",
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}

	return c
}
