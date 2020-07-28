package domain

import (
	"errors"
	"time"
)

// VoyageNumber uniquely identifies a particular Voyage.
type VoyageNumber string

// Voyage is a uniquely identifiable series of carrier movements.
type Voyage struct {
	Number   VoyageNumber
	Schedule Schedule
}

// New creates a voyage with a voyage number and a provided schedule.
func New(n VoyageNumber, s Schedule) *Voyage {
	return &Voyage{Number: n, Schedule: s}
}

// Schedule describes a voyage schedule.
type Schedule struct {
	CarrierMovements []CarrierMovement
}

// CarrierMovement is a vessel voyage from one location to another.
type CarrierMovement struct {
	DepartureLocation UNLocode
	ArrivalLocation   UNLocode
	DepartureTime     time.Time
	ArrivalTime       time.Time
}

// ErrUnknown is used when a voyage could not be found.
var ErrUnknown = errors.New("unknown voyage")

// Repository provides access a voyage store.
type Repository interface {
	Find(VoyageNumber) (*Voyage, error)
}
