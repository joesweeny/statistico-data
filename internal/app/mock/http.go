package mock

import (
	"net/http"
)

func HttpClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (r roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}
