package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/factory"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FixtureHandler struct {
	fixtureRepo app.FixtureRepository
	factory     *factory.FixtureFactory
	logger 		*logrus.Logger
}

func RoutePath(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Welcome to the Statistico Data API")
}

func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Healthcheck OK")
}

func (f FixtureHandler) SeasonFixtures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	after, err := parseDateQuery(r.URL.Query(), "date_after")

	if err == errTimeParse {
		failResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	before, err := parseDateQuery(r.URL.Query(), "date_before")

	if err == errTimeParse {
		failResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	fixtures, err := f.fixtureRepo.BySeasonIDBetween(uint64(id), after, before)

	if err != nil {
		f.logger.Errorf("Error fetching fixtures. Url %s, error %s", r.URL.RawPath, err.Error())
		errorResponse(w, 500, err)
		return
	}

	var response []*Fixture

	for _, fix := range fixtures {
		x, err := f.factory.BuildRestFixture(&fix)

		if err != nil {
			f.logger.Errorf("Error hydrating fixture. Url %s, error %s", r.URL.RawPath, err.Error())
			errorResponse(w, 500, err)
			return
		}

		response = append(response, x)
	}

	successResponse(w, response)
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

func NewFixtureHandler(r app.FixtureRepository, f *factory.FixtureFactory, l *logrus.Logger) *FixtureHandler {
	return &FixtureHandler{fixtureRepo: r, factory: f, logger: l}
}
