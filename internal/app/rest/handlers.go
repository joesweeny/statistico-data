package rest

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

var errTimeParse = errors.New("date provided in request is not a valid RFC3339 formatted date")

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

	after := query.Get("date_after")
	before := query.Get("date_before")

	if after != "" {
		_, err := time.Parse(time.RFC3339, after)

		if err != nil {
			jsendFailResponse(w, http.StatusUnprocessableEntity, errTimeParse)
			return
		}
	}

	if after != "" {
		_, err := time.Parse(time.RFC3339, before)

		if err != nil {
			jsendFailResponse(w, http.StatusUnprocessableEntity, errTimeParse)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Season %s \n", id)
	_, _ = fmt.Fprintf(w, "After date %s", after)
	_, _ = fmt.Fprintf(w, "Before date %s", before)
}