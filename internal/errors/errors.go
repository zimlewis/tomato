package errors

import "errors"

var (
	ErrDidNotStart = errors.New("Did not start the tomato session")
	ErrBadgerDB = errors.New("Cannot interact with the DB")
)
