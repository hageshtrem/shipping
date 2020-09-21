package application

import "time"

// TrackingID uniquely identifies a particular cargo.
type TrackingID string

// Cargo is a read model for tracking views.
type Cargo struct {
	TrackingID           string    `json:"tracking_id"`
	StatusText           string    `json:"status_text"`
	Origin               string    `json:"origin"`
	Destination          string    `json:"destination"`
	ETA                  time.Time `json:"eta"`
	NextExpectedActivity string    `json:"next_expected_activity"`
	ArrivalDeadline      time.Time `json:"arrival_deadline"`
	// Events               []Event   `json:"events"`
}

type CargoViewModelRepository interface {
	Store(cargo *Cargo) error
	Find(id string) (*Cargo, error)
	FindAll() []*Cargo
}
