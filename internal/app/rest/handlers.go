package rest

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"time"
)

func routePath(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Welcome to the Statistico Data API")
}

func healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Healthcheck OK")
}

func seasonFixtures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	query := r.URL.Query()

	after, err := parseDateQuery(query, "date_after")

	if err == errTimeParse {
		failResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	before, err := parseDateQuery(query, "date_before")

	if err == errTimeParse {
		failResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Season %s \n", id)

	if after != nil {
		_, _ = fmt.Fprintf(w, "After date %s", after)
	}

	if before != nil {
		_, _ = fmt.Fprintf(w, "Before date %s", before)
	}
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
