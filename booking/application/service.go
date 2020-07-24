package application

import (
	"booking/domain"
	"errors"
	"time"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// Service is the interface that provides booking methods.
type Service interface {
	// BookNewCargo registers a new cargo in the tracking system, not yet
	// routed.
	BookNewCargo(origin domain.UNLocode, destination domain.UNLocode, deadline time.Time) (domain.TrackingID, error)
}

// NewService creates a booking service with necessary dependencies.
func NewService(cargos domain.CargoRepository) Service {
	return &service{
		cargos: cargos,
	}
}

type service struct {
	cargos domain.CargoRepository
}

func (s *service) BookNewCargo(origin domain.UNLocode, destination domain.UNLocode, deadline time.Time) (domain.TrackingID, error) {
	if origin == "" || destination == "" || deadline.IsZero() {
		return "", ErrInvalidArgument
	}

	id := s.cargos.NextTrackingID()
	rs := domain.RouteSpecification{
		Origin:          origin,
		Destination:     destination,
		ArrivalDeadline: deadline,
	}

	c := domain.NewCargo(id, rs)

	if err := s.cargos.Store(c); err != nil {
		return "", err
	}

	return c.TrackingID, nil
}
