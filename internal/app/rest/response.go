package rest

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func successResponse(w http.ResponseWriter, status int, payload interface{}) {
	response := response{
		Message: "success",
		Data:    payload,
	}

	jsonResponse(w, status, response)
}

func failResponse(w http.ResponseWriter, status int, error error) {
	response := response{
		Message: "fail",
		Data: []errorMessage{
			{
				Message: error.Error(),
				Code:    1,
			},
		},
	}

	jsonResponse(w, status, response)
}

func errorResponse(w http.ResponseWriter, status int, error error) {
	response := response{
		Message: "error",
		Data: []errorMessage{
			{
				Message: error.Error(),
				Code:    1,
			},
		},
	}

	jsonResponse(w, status, response)
}
