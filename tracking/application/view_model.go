package application

import "time"

// TrackingID uniquely identifies a particular cargo.
type TrackingID string

// Cargo is a read model for tracking views.
type Cargo struct {
	TrackingID           string
	StatusText           string
	Origin               string
	Destination          string
	ETA                  time.Time
	NextExpectedActivity string
	ArrivalDeadline      time.Time
	Events               []Event
}

// Event is a read model for tracking views.
type Event struct {
	Description string
	Expected    bool
}

// CargoViewModelRepository provides access a cargo view model store.
type CargoViewModelRepository interface {
	Store(cargo *Cargo) error
	Find(id string) (*Cargo, error)
	FindAll() []*Cargo
}
