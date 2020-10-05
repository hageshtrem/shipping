package application

import (
	"log"
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
