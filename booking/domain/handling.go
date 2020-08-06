package domain

// TODO: It would make sense not having the cargo package depend on handling.

import (
	"errors"
)

// HandlingActivity represents how and where a cargo can be handled, and can
// be used to express predictions about what is expected to happen to a cargo
// in the future.
type HandlingActivity struct {
	Type         HandlingEventType
	Location     UNLocode
	VoyageNumber VoyageNumber
}

// HandlingEvent is used to register the event when, for instance, a cargo is
// unloaded from a carrier at a some location at a given time.
type HandlingEvent struct {
	TrackingID TrackingID
	Activity   HandlingActivity
}

// HandlingEventType describes type of a handling event.
type HandlingEventType int

// Valid handling event types.
const (
	NotHandled HandlingEventType = iota
	Load
	Unload
	Receive
	Claim
	Customs
)

func (t HandlingEventType) String() string {
	switch t {
	case NotHandled:
		return "Not Handled"
	case Load:
		return "Load"
	case Unload:
		return "Unload"
	case Receive:
		return "Receive"
	case Claim:
		return "Claim"
	case Customs:
		return "Customs"
	}

	return ""
}

// HandlingHistory is the handling history of a cargo.
type HandlingHistory struct {
	HandlingEvents []HandlingEvent
}

// MostRecentlyCompletedEvent returns most recently completed handling event.
func (h HandlingHistory) MostRecentlyCompletedEvent() (HandlingEvent, error) {
	if len(h.HandlingEvents) == 0 {
		return HandlingEvent{}, errors.New("delivery history is empty")
	}

	return h.HandlingEvents[len(h.HandlingEvents)-1], nil
}

// HandlingEventRepository provides access a handling event store.
type HandlingEventRepository interface {
	Store(e HandlingEvent)
	QueryHandlingHistory(TrackingID) HandlingHistory
}
