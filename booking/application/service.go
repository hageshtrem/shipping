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

	// RequestPossibleRoutesForCargo requests a list of itineraries describing
	// possible routes for this cargo.
	RequestPossibleRoutesForCargo(id domain.TrackingID) []domain.Itinerary
}

// NewService creates a booking service with necessary dependencies.
func NewService(cargos domain.CargoRepository, rs domain.RoutingService) Service {
	return &service{
		cargos:         cargos,
		routingService: rs,
	}
}

type service struct {
	cargos         domain.CargoRepository
	routingService domain.RoutingService
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

func (s *service) RequestPossibleRoutesForCargo(id domain.TrackingID) []domain.Itinerary {
	if id == "" {
		return nil
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return []domain.Itinerary{}
	}

	return s.routingService.FetchRoutesForSpecification(c.RouteSpecification)
}
