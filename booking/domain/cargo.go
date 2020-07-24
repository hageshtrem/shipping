package domain

import(
	"time"
	"errors"
)

// TrackingID uniquely identifies a particular cargo.
type TrackingID string

// Cargo is the central class in the domain model.
type Cargo struct {
	TrackingID         TrackingID
	Origin             UNLocode
	RouteSpecification RouteSpecification
	// Itinerary          Itinerary
	// Delivery           Delivery
}

// NewCargo creates a new, unrouted cargo.
func NewCargo(id TrackingID, rs RouteSpecification) *Cargo {
	// itinerary := Itinerary{}
	// history := HandlingHistory{make([]HandlingEvent, 0)}

	return &Cargo{
		TrackingID:         id,
		Origin:             rs.Origin,
		RouteSpecification: rs,
		// Delivery:           DeriveDeliveryFrom(rs, itinerary, history),
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