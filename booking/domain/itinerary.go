package domain

import (
	"time"
)

// Leg describes the transportation between two locations on a voyage.
type Leg struct {
	VoyageNumber   VoyageNumber `json:"voyage_number"`
	LoadLocation   UNLocode     `json:"from"`
	UnloadLocation UNLocode     `json:"to"`
	LoadTime       time.Time    `json:"load_time"`
	UnloadTime     time.Time    `json:"unload_time"`
}

// NewLeg creates a new itinerary leg.
func NewLeg(voyageNumber VoyageNumber, loadLocation, unloadLocation UNLocode, loadTime, unloadTime time.Time) Leg {
	return Leg{
		VoyageNumber:   voyageNumber,
		LoadLocation:   loadLocation,
		UnloadLocation: unloadLocation,
		LoadTime:       loadTime,
		UnloadTime:     unloadTime,
	}
}

// Itinerary specifies steps required to transport a cargo from its origin to
// destination.
type Itinerary struct {
	Legs []Leg `json:"legs"`
}
