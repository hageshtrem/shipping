package application

import "errors"

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrUnknownCargo is used when a cargo could not be found.
var ErrUnknownCargo = errors.New("unknown cargo")

// Service is the interface that provides the basic Track method.
type Service interface {
	// Track returns a cargo matching a tracking ID.
	Track(id string) (Cargo, error)
}

type service struct {
	cargos CargoViewModelRepository
}

func (s *service) Track(id string) (Cargo, error) {
	if id == "" {
		return Cargo{}, ErrInvalidArgument
	}
	c, err := s.cargos.Find(id)
	if err != nil {
		return Cargo{}, err
	}
	return *c, nil
}

// NewService returns a new instance of the default Service.
func NewService(cargos CargoViewModelRepository) Service {
	return &service{cargos}
}
