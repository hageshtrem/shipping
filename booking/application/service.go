package application

// TODO: add transactions
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

	// LoadCargo returns a read model of a shipping.
	LoadCargo(id domain.TrackingID) (Cargo, error)

	// Cargos returns a list of all cargos that have been booked.
	Cargos() []Cargo
}

// NewService creates a booking service with necessary dependencies.
func NewService(cargos domain.CargoRepository, locations domain.LocationRepository, rs domain.RoutingService, eventService EventService) Service {
	return &service{
		cargos:         cargos,
		locations:      locations,
		routingService: rs,
		eventService:   eventService,
	}
}

type service struct {
	cargos         domain.CargoRepository
	locations      domain.LocationRepository
	routingService domain.RoutingService
	eventService   EventService
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

	if err := s.eventService.NewCargoBooked(c); err != nil {
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

	if err := s.cargos.Store(c); err != nil {
		return err
	}

	return s.eventService.CargoToRouteAssigned(c)
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

	return s.eventService.DestinationChanged(c)
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

func (s *service) LoadCargo(id domain.TrackingID) (Cargo, error) {
	if id == "" {
		return Cargo{}, ErrInvalidArgument
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return Cargo{}, err
	}

	return assemble(c), nil
}

func (s *service) Cargos() []Cargo {
	var result []Cargo
	for _, c := range s.cargos.FindAll() {
		result = append(result, assemble(c))
	}
	return result
}

// Cargo is a read model for booking views.
type Cargo struct {
	ArrivalDeadline time.Time    `json:"arrival_deadline"`
	Destination     string       `json:"destination"`
	Legs            []domain.Leg `json:"legs,omitempty"`
	Misrouted       bool         `json:"misrouted"`
	Origin          string       `json:"origin"`
	Routed          bool         `json:"routed"`
	TrackingID      string       `json:"tracking_id"`
}

func assemble(c *domain.Cargo) Cargo {
	return Cargo{
		TrackingID:      string(c.TrackingID),
		Origin:          string(c.Origin),
		Destination:     string(c.RouteSpecification.Destination),
		Misrouted:       c.Delivery.RoutingStatus == domain.Misrouted,
		Routed:          !c.Itinerary.IsEmpty(),
		ArrivalDeadline: c.RouteSpecification.ArrivalDeadline,
		Legs:            c.Itinerary.Legs,
	}
}
