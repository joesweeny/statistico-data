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

var c = model.Country{
	uuid.UUID{},
	1,
	"England",
	"Europe",
	"ENG",
	time.Now(),
	time.Now(),
}

func TestInsert(t *testing.T) {
	conn := getConnection()

	t.Run("insert increases table count", func(t *testing.T) {
		defer conn.Close()
		tearDown(conn)
		repo := NewPostgresCountryRepository(conn)
		err := repo.Insert(c)

		if err != nil {
			t.Error("Error when inserting record into the database")
		}

		row := conn.QueryRow("select count(*) from country")

		if err != nil {
			panic(err)
		}

		var count int
		err = row.Scan(&count)

		if err != nil {
			panic(err)
		}

		fmt.Printf("count: %d\n", count)
	})
}

func getConnection() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	return db
}

func tearDown(db *sql.DB) {
	db.Exec("delete from country")
}
