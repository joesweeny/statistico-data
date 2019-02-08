package squad

import (
	"database/sql"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/joesweeny/statshub/internal/season"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresSquadRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newSquad(i, i+1)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_squad")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})
}

func TestBySeasonAndTeam(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresSquadRepository{Connection: conn}

	t.Run("squad can be retrieved by season and team IDs", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newSquad(45, 986)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.BySeasonAndTeam(45, 986)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(45, r.SeasonID)
		a.Equal(986, m.TeamID)
		a.Equal([]int{34, 57, 89}, m.PlayerIDs)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.BySeasonAndTeam(99, 2435870)

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
	repo := PostgresSquadRepository{Connection: conn}

	t.Run("modifies existing squad", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newSquad(25, 62)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		m.PlayerIDs = []int{432, 567, 2, 87095}
		m.UpdatedAt = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		if err := repo.Update(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.BySeasonAndTeam(25, 62)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(25, r.SeasonID)
		a.Equal(62, m.TeamID)
		a.Equal([]int{432, 567, 2, 87095}, m.PlayerIDs)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if player does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newSquad(146, 200)

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

func TestAll(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresSquadRepository{Connection: conn}

	t.Run("returns all squad records from the database", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		squads := []model.Squad{
			*newSquad(39, 25067),
			*newSquad(99, 98),
			*newSquad(301, 2),
			*newSquad(23, 6),
			*newSquad(39, 1902),
		}

		for _, squad := range squads {
			repo.Insert(&squad)
		}

		all, err := repo.All()

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		a := assert.New(t)

		a.Equal(5, len(all))
		a.Equal(all[:1], squads[3:4])
		a.Equal(all[1:2], squads[4:5])
		a.Equal(all[2:3], squads[0:1])
		a.Equal(all[3:4], squads[1:2])
		a.Equal(all[4:], squads[2:3])
	})
}

func TestCurrentSeason(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresSquadRepository{Connection: conn}
	seasonRepo := season.PostgresSeasonRepository{Connection: conn}

	t.Run("returns squads only for current season", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		seasonRepo.Insert(newSeason(true, 39))
		seasonRepo.Insert(newSeason(false, 99))
		seasonRepo.Insert(newSeason(false, 4502))
		seasonRepo.Insert(newSeason(true, 23))

		squads := []model.Squad{
			*newSquad(39, 25067),
			*newSquad(99, 98),
			*newSquad(301, 2),
			*newSquad(23, 6),
			*newSquad(39, 1902),
		}

		for _, squad := range squads {
			repo.Insert(&squad)
		}

		current, err := repo.CurrentSeason()

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		a := assert.New(t)

		a.Equal(3, len(current))
		a.Equal(current[:1], squads[3:4])
		a.Equal(current[1:2], squads[4:])
		a.Equal(current[2:3], squads[0:1])
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
		_, err := db.Exec("delete from sportmonks_squad")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newSquad(season, team int) *model.Squad {
	return &model.Squad{
		SeasonID:  season,
		TeamID:    team,
		PlayerIDs: []int{34, 57, 89},
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

func newSeason(current bool, seasonId int) *model.Season {
	return &model.Season{
		ID:        seasonId,
		Name:      "Current Season",
		LeagueID:  4,
		IsCurrent: current,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}
