package errors

import "errors"

var ErrorDatabaseConnection = errors.New("error connecting the to the database")
var ErrorNotFound = errors.New("the resource requested does not exist")
