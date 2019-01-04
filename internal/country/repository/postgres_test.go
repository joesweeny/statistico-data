package repository

import (
	"testing"
	"os"
	"fmt"
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
	"time"
	_ "github.com/lib/pq"
)

var host = os.Getenv("DB_HOST")
var port = os.Getenv("DB_PORT")
var user = os.Getenv("DB_USER")
var pass = os.Getenv("DB_PASSWORD")
var name = os.Getenv("DB_NAME")

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := NewPostgresCountryRepository(conn)

	t.Run("insert increases table count", func(t *testing.T) {
		t.Helper()
		for i := 1; i < 4; i++ {
			var id = uuid.Must(uuid.NewV4())
			c := newCountry(id)

			if err := repo.Insert(c); err != nil {
				t.Fatalf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from country")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Fatalf("Error when scanning rows returned by the database: %s", err.Error())
			}

			if i != count {
				t.Fatalf("Expected %d, got %d", i, count)
			}
		}

		cleanUp()
	})

	t.Run("insert returns error when ID primary key violates unique constraint", func(t *testing.T) {
		t.Helper()
		c := newCountry(uuid.UUID{})

		if err := repo.Insert(c); err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}

		cleanUp()
	})

	conn.Close()
}

func getConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", psqlInfo)

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
		u,
		1,
		"England",
		"Europe",
		"ENG",
		time.Now(),
		time.Now(),
	}

	return c;
}
