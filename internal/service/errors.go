package service

import "errors"

// ErrNotFound is returned when the requested entity does not exist.
var ErrNotFound = errors.New("entity not found")
