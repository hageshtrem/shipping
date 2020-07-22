package domain

import(
	"time"
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

// RouteSpecification Contains information about a route: its origin,
// destination and arrival deadline.
type RouteSpecification struct {
	Origin          UNLocode
	Destination     UNLocode
	ArrivalDeadline time.Time
}