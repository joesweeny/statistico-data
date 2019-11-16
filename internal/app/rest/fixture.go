package rest

import (
	"github.com/julienschmidt/httprouter"
	"github.com/statistico/statistico-data/internal/app"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FixtureHandler struct {
	fixtureRepo app.FixtureRepository
	factory     *FixtureFactory
}

func (f FixtureHandler) SeasonFixtures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query, err := parseFixtureQuery(r, ps)

	if err == errBadRequest {
		failResponse(w, http.StatusBadRequest, err)
		return
	}

	if err == errTimeParse {
		failResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	fixtures, err := f.fixtureRepo.Get(query)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, internalServerError)
		return
	}

	response := fixtureResponse{Fixtures: []Fixture{}}

	for _, fix := range fixtures {
		f, err := f.factory.BuildFixture(&fix)

		if err != nil {
			errorResponse(w, http.StatusInternalServerError, internalServerError)
			return
		}

		response.Fixtures = append(response.Fixtures, *f)
	}

	successResponse(w, http.StatusOK, response)
}

func parseFixtureQuery(r *http.Request, ps httprouter.Params) (app.FixtureRepositoryQuery, error) {
	query := app.FixtureRepositoryQuery{}

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		return query, errBadRequest
	}

	from, err := parseDateQuery(r.URL.Query(), "date_from")

	if err == errTimeParse {
		return query, err
	}

	to, err := parseDateQuery(r.URL.Query(), "date_to")

	if err == errTimeParse {
		return query, err
	}

	var sort string

	if r.URL.Query().Get("sort") != "" {
		sort = r.URL.Query().Get("sort")
	}

	seasonID := uint64(id)
	query.SeasonID = &seasonID
	query.DateFrom = from
	query.DateTo = to
	query.SortBy = &sort

	return query, nil
}

func parseDateQuery(query url.Values, key string) (*time.Time, error) {
	val := query.Get(key)

	if val == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, val)

	if err != nil {
		return nil, errTimeParse
	}

	return &t, nil
}

func NewFixtureHandler(r app.FixtureRepository, f *FixtureFactory) *FixtureHandler {
	return &FixtureHandler{fixtureRepo: r, factory: f}
}
