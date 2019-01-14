package country

import (
	"github.com/stretchr/testify/mock"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/satori/go.uuid"
	"github.com/jonboulle/clockwork"
	"testing"
	"github.com/pkg/errors"
	"time"
)

type mockRepository struct {
	mock.Mock
}

func (r *mockRepository) Insert(c model.Country) error {
	args := r.Called(c)
	return args.Error(0)
}

func (r *mockRepository) Update(c model.Country) error {
	args := r.Called(c)
	return args.Error(0)
}

func (r *mockRepository) GetById(id uuid.UUID) (model.Country, error) {
	args := r.Called(id)
	return args.Get(0).(model.Country), args.Error(0)
}

func (r *mockRepository) GetByExternalId(id int) (model.Country, error) {
	args := r.Called(id)
	return args.Get(0).(model.Country), args.Error(1)
}

func TestHandle(t *testing.T) {
	c := country()

	t.Run("inserts new record into the database if does not exist", func (t *testing.T) {
		t.Helper()

		var repo = new(mockRepository)
		var factory = Factory{clockwork.NewFakeClock()}
		var handler = Handler{repo, factory}
		var e = errors.New("Not found")

		repo.On("GetByExternalId", 180).Return(model.Country{}, e)
		repo.On("Insert", mock.Anything).Return(nil)

		repo.AssertNotCalled(t, "Update", mock.Anything)

		if err := handler.Handle(c); err != nil {
			t.Fatalf("Expected return value to be nil, got %s", err)
		}
	})

	t.Run("updates record if it already exists", func (t *testing.T) {
		t.Helper()

		var repo = new(mockRepository)
		var factory = Factory{clockwork.NewFakeClock()}
		var handler = Handler{repo, factory}

		fetched := model.Country{
			ID:         uuid.NewV4(),
			ExternalID: 180,
			Name:       "United Kingdom",
			Continent:  "Europe",
			ISO:        "ENG",
			CreatedAt:  time.Unix(1546965200, 0),
			UpdatedAt:  time.Unix(1546965200, 0),
		}

		repo.On("GetByExternalId", 180).Return(fetched, nil)
		repo.On("Update", mock.Anything).Return(nil)

		repo.AssertNotCalled(t, "Insert", mock.Anything)

		if err := handler.Handle(c); err != nil {
			t.Fatalf("Expected return value to be nil, got %s", err)
		}
	})

	t.Run("error is returned if unable to insert into repository", func (t *testing.T) {
		t.Helper()

		var repo = new(mockRepository)
		var factory = Factory{clockwork.NewFakeClock()}
		var handler = Handler{repo, factory}
		var e = errors.New("Not found")

		repo.On("GetByExternalId", 180).Return(model.Country{}, e)
		repo.On("Insert", mock.Anything).Return(errors.New("Unable to insert"))

		repo.AssertNotCalled(t, "Update", mock.Anything)

		if err := handler.Handle(c); err == nil {
			t.Fatalf("Expected return value to be %s, got nil", err)
		}
	})
}

func country() sportmonks.Country {
	country := sportmonks.Country{
		ID: 180,
		Name: "England",
		Extra: struct {
			Continent string `json:"continent"`
			SubRegion string `json:"sub_region"`
			WorldRegion string `json:"world_region"`
			Fifa interface{} `json:"fifa,string"`
			ISO string `json:"iso"`
			Longitude string `json:"longitude"`
			Latitude string `json:"latitude"`
		} {
			Continent:   "Europe",
			SubRegion:   "Western Europe",
			WorldRegion: "Europe",
			Fifa:        "ENG",
			ISO:         "ENG",
		},
	}

	return country
}