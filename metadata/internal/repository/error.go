package repository

import "errors"

// ErrNotFound is returned when requested record in not found.
var ErrNotFound = errors.New("not found")
