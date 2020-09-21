package application

import "errors"

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")
var ErrUnknownCargo = errors.New("unknown cargo")

// Service is the interface that provides the basic Track method.
type Service interface {
	// Track returns a cargo matching a tracking ID.
	Track(id string) (Cargo, error)
}
