package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	sea := ps.ByName("id")

	query := r.URL.Query()

	after := query.Get("date_after")
	before := query.Get("date_before")

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Season %s \n", sea)

	if after != "" {
		_, _ = fmt.Fprintf(w, "After date %s", after)
	}

	if before != "" {
		_, _ = fmt.Fprintf(w, "Before date %s", before)
	}
}
