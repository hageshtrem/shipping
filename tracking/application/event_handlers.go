package application

import (
	"fmt"
	"log"
	"strings"
	"time"
	booking "tracking/pb/booking/pb"

	"google.golang.org/protobuf/proto"
)

// EventHandler is is the interface that provides events (RabbitMQ) handling.
type EventHandler interface {
	Handle(event proto.Message) error
}

type newCargoBookedEventHandler struct {
	cargos CargoViewModelRepository
}

// NewCargoBookedEventHandler creates an event handler for CargoBooked event.
func NewCargoBookedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &newCargoBookedEventHandler{cargos}
}

func (eh *newCargoBookedEventHandler) Handle(event proto.Message) error {
	newCargo := event.(*booking.NewCargoBooked)
	log.Printf("New Cargo booked: %v", newCargo)
	cargo := Cargo{
		TrackingID:      newCargo.GetTrackingId(),
		Origin:          newCargo.GetOrigin(),
		Destination:     newCargo.GetDestination(),
		ArrivalDeadline: newCargo.GetArrivalDeadline().AsTime(),
	}
	return eh.cargos.Store(&cargo)
}

type cargoDestinationChangedEventHandler struct {
	cargos CargoViewModelRepository
}

// NewCargoDestinationChangedEventHandler creates an event handler for CargoDestinationChanged event.
func NewCargoDestinationChangedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &cargoDestinationChangedEventHandler{cargos}
}

func (eh *cargoDestinationChangedEventHandler) Handle(event proto.Message) error {
	e := event.(*booking.CargoDestinationChanged)
	log.Printf("Cargo %s destination changed %s", e.TrackingId, e.Destination)
	c, err := eh.cargos.Find(e.GetTrackingId())
	if err != nil {
		return err
	}

	c.Destination = e.Destination
	return eh.cargos.Store(c)
}

type cargoToRouteAssignedEventHandler struct {
	cargos CargoViewModelRepository
}

// NewCargoToRouteAssignedEventHandler creates an event handler for CargoToRouteAssigned event.
func NewCargoToRouteAssignedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &cargoToRouteAssignedEventHandler{cargos}
}

func (eh *cargoToRouteAssignedEventHandler) Handle(event proto.Message) error {
	e := event.(*booking.CargoToRouteAssigned)
	log.Printf("Cargo %s assigned to the route", e.GetTrackingId())
	c, err := eh.cargos.Find(e.GetTrackingId())
	if err != nil {
		return err
	}

	c.ETA = e.GetEta().AsTime()
	return eh.cargos.Store(c)
}

type cargoWasHandledEventHandler struct {
	cargos CargoViewModelRepository
}

// NewCargoWasHandledEventHandler creates an event handler for CargoWasHandled event.
func NewCargoWasHandledEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &cargoWasHandledEventHandler{cargos}
}

func (eh *cargoWasHandledEventHandler) Handle(event proto.Message) error {
	e := event.(*booking.CargoWasHandled)
	log.Printf("Cargo was handled {%v}", e.Delivery)
	c, err := eh.cargos.Find(e.GetTrackingId())
	if err != nil {
		return err
	}

	delivery := e.GetDelivery()

	c.StatusText = delivery.GetTransportStatus()
	c.NextExpectedActivity = nextExpectedActivity(delivery.GetNextExpectedActivity())
	c.ETA = delivery.GetEta().AsTime()
	c.Events = append(c.Events, assembleEvent(delivery.GetLastEvent()))

	return eh.cargos.Store(c)
}

func nextExpectedActivity(activity *booking.HandlingActivity) string {
	prefix := "Next expected activity is to"

	switch activity.Type {
	case booking.HandlingEventType_Load:
		return fmt.Sprintf("%s %s cargo onto voyage %s in %s.", prefix, strings.ToLower(activity.Type.String()), activity.VoyageNumber, activity.Location)
	case booking.HandlingEventType_Unload:
		return fmt.Sprintf("%s %s cargo off of voyage %s in %s.", prefix, strings.ToLower(activity.Type.String()), activity.VoyageNumber, activity.Location)
	case booking.HandlingEventType_NotHandled:
		return "There are currently no expected activities for this shipping."
	}

	return fmt.Sprintf("%s %s cargo in %s.", prefix, strings.ToLower(activity.Type.String()), activity.Location)
}

func assembleEvent(e *booking.HandlingEvent) Event {
	var description string

	switch e.Activity.Type {
	case booking.HandlingEventType_NotHandled:
		description = "Cargo has not yet been received."
	case booking.HandlingEventType_Receive:
		description = fmt.Sprintf("Received in %s, at %s", e.Activity.Location, time.Now().Format(time.RFC3339))
	case booking.HandlingEventType_Load:
		description = fmt.Sprintf("Loaded onto voyage %s in %s, at %s.", e.Activity.VoyageNumber, e.Activity.Location, time.Now().Format(time.RFC3339))
	case booking.HandlingEventType_Unload:
		description = fmt.Sprintf("Unloaded off voyage %s in %s, at %s.", e.Activity.VoyageNumber, e.Activity.Location, time.Now().Format(time.RFC3339))
	case booking.HandlingEventType_Claim:
		description = fmt.Sprintf("Claimed in %s, at %s.", e.Activity.Location, time.Now().Format(time.RFC3339))
	case booking.HandlingEventType_Customs:
		description = fmt.Sprintf("Cleared customs in %s, at %s.", e.Activity.Location, time.Now().Format(time.RFC3339))
	default:
		description = "[Unknown status]"
	}

	return Event{
		Description: description,
		Expected:    false, // TODO
	}
}
