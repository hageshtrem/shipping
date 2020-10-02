package application

import (
	"log"
	booking "tracking/pb/booking/pb"

	"google.golang.org/protobuf/proto"
)

type EventHandler interface {
	Handle(event proto.Message) error
}

type newCargoBookedEventHandler struct {
	cargos CargoViewModelRepository
}

func NewCargoBookedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &newCargoBookedEventHandler{cargos}
}

func (eh *newCargoBookedEventHandler) Handle(event proto.Message) error {
	newCargo := event.(*booking.NewCargoBooked)
	log.Printf("New Cargo booked: %v", newCargo)
	cargo := Cargo{
		TrackingID:      newCargo.TrackingId,
		Origin:          newCargo.Origin,
		Destination:     newCargo.Destination,
		ArrivalDeadline: newCargo.ArrivalDeadline.AsTime(),
	}
	return eh.cargos.Store(&cargo)
}

type cargoDestinationChangedEventHandler struct {
	cargos CargoViewModelRepository
}

func NewCargoDestinationChangedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &cargoDestinationChangedEventHandler{cargos}
}

func (eh *cargoDestinationChangedEventHandler) Handle(event proto.Message) error {
	e := event.(*booking.CargoDestinationChanged)
	log.Printf("Cargo %s destination changed %s", e.TrackingId, e.Destination)
	c, err := eh.cargos.Find(e.TrackingId)
	if err != nil {
		return err
	}

	c.Destination = e.Destination
	return eh.cargos.Store(c)
}
