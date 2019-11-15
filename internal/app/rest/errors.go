package rest

import "errors"

var errTimeParse = errors.New("date provided in request is not a valid RFC3339 formatted date")
var errBadRequest = errors.New("request provided is not in a valid format")
var internalServerError = errors.New("internal server error")

