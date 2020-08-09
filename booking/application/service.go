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

	// ChangeDestination changes the destination of a shipping.
	ChangeDestination(id domain.TrackingID, destination domain.UNLocode) error

	// Locations returns a list of registered locations.
	Locations() []Location
}

// NewService creates a booking service with necessary dependencies.
func NewService(cargos domain.CargoRepository, locations domain.LocationRepository, rs domain.RoutingService) Service {
	return &service{
		cargos:         cargos,
		locations:      locations,
		routingService: rs,
	}
}

type service struct {
	cargos         domain.CargoRepository
	locations      domain.LocationRepository
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

func (s *service) ChangeDestination(id domain.TrackingID, destination domain.UNLocode) error {
	if id == "" || destination == "" {
		return ErrInvalidArgument
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return err
	}

	l, err := s.locations.Find(destination)
	if err != nil {
		return err
	}

	c.SpecifyNewRoute(domain.RouteSpecification{
		Origin:          c.Origin,
		Destination:     l.UNLocode,
		ArrivalDeadline: c.RouteSpecification.ArrivalDeadline,
	})

	if err := s.cargos.Store(c); err != nil {
		return err
	}

	return nil
}

func (s *service) Locations() []Location {
	var result []Location
	for _, v := range s.locations.FindAll() {
		result = append(result, Location{
			UNLocode: string(v.UNLocode),
			Name:     v.Name,
		})
	}
	return result
}

// Location is a read model for booking views.
type Location struct {
	UNLocode string `json:"locode"`
	Name     string `json:"name"`
}
