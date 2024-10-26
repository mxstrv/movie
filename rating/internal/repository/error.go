package repository

import "errors"

// ErrNotFound is returned when requested rating is not found.
var ErrNotFound = errors.New("not found")
