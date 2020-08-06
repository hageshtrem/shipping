package domain

import (
	"errors"
	"time"
)

// TrackingID uniquely identifies a particular cargo.
type TrackingID string

// Cargo is the central class in the domain model.
type Cargo struct {
	TrackingID         TrackingID
	Origin             UNLocode
	RouteSpecification RouteSpecification
	Itinerary          Itinerary
	Delivery           Delivery
}

// NewCargo creates a new, unrouted cargo.
func NewCargo(id TrackingID, rs RouteSpecification) *Cargo {
	itinerary := Itinerary{}
	history := HandlingHistory{make([]HandlingEvent, 0)}

	return &Cargo{
		TrackingID:         id,
		Origin:             rs.Origin,
		RouteSpecification: rs,
		Delivery:           DeriveDeliveryFrom(rs, itinerary, history),
	}
}

// CargoRepository provides access a cargo store.
type CargoRepository interface {
	Store(cargo *Cargo) error
	Find(id TrackingID) (*Cargo, error)
	FindAll() []*Cargo
	NextTrackingID() TrackingID
}

// ErrUnknownCargo is used when a cargo could not be found.
var ErrUnknownCargo = errors.New("unknown cargo")

// RouteSpecification Contains information about a route: its origin,
// destination and arrival deadline.
type RouteSpecification struct {
	Origin          UNLocode
	Destination     UNLocode
	ArrivalDeadline time.Time
}

// IsSatisfiedBy checks whether provided itinerary satisfies this
// specification.
func (s RouteSpecification) IsSatisfiedBy(itinerary Itinerary) bool {
	return itinerary.Legs != nil &&
		s.Origin == itinerary.InitialDepartureLocation() &&
		s.Destination == itinerary.FinalArrivalLocation()
}

// RoutingStatus describes status of cargo routing.
type RoutingStatus int

// Valid routing statuses.
const (
	NotRouted RoutingStatus = iota
	Misrouted
	Routed
)

func (s RoutingStatus) String() string {
	switch s {
	case NotRouted:
		return "Not routed"
	case Misrouted:
		return "Misrouted"
	case Routed:
		return "Routed"
	}
	return ""
}

// TransportStatus describes status of cargo transportation.
type TransportStatus int

// Valid transport statuses.
const (
	NotReceived TransportStatus = iota
	InPort
	OnboardCarrier
	Claimed
	Unknown
)

func (s TransportStatus) String() string {
	switch s {
	case NotReceived:
		return "Not received"
	case InPort:
		return "In port"
	case OnboardCarrier:
		return "Onboard carrier"
	case Claimed:
		return "Claimed"
	case Unknown:
		return "Unknown"
	}
	return ""
}
