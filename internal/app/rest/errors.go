package rest

import "errors"

var errTimeParse = errors.New("date provided in request is not a valid RFC3339 formatted date")
