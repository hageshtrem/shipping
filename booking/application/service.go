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

	// AssignCargoToRoute assigns a cargo to the route specified by the
	// itinerary.
	AssignCargoToRoute(id domain.TrackingID, itinerary domain.Itinerary) error
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

func (s *service) AssignCargoToRoute(id domain.TrackingID, itinerary domain.Itinerary) error {
	if id == "" || len(itinerary.Legs) == 0 {
		return ErrInvalidArgument
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return err
	}

	c.AssignToRoute(itinerary)

	return s.cargos.Store(c)
}
