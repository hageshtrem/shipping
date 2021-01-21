package application

import (
	"booking/domain"
	handling "booking/pb/handling/pb"

	"google.golang.org/protobuf/proto"
)

// EventHandler is an abstraction for an event handler.
type EventHandler interface {
	Handle(event proto.Message) error
}

type cargoHandledEventHandler struct {
	cargos       domain.CargoRepository
	eventService EventService
}

// NewCargoHandledEventHandler return a handler for the CargoHandled event.
func NewCargoHandledEventHandler(cargos domain.CargoRepository, es EventService) EventHandler {
	return &cargoHandledEventHandler{cargos, es}
}

func (eh *cargoHandledEventHandler) Handle(event proto.Message) error {
	handlingEvent := event.(*handling.HandlingEvent)
	domainHandlingEvent := domain.HandlingEvent{
		TrackingID: domain.TrackingID(handlingEvent.TrackingId),
		Activity:   decodeHandlingActivity(handlingEvent.Activity),
	}
	cargo, err := eh.cargos.Find(domain.TrackingID(handlingEvent.TrackingId))
	if err != nil {
		return err
	}

	// TODO: It needs to be transactional
	cargo.DeriveDeliveryProgress(domainHandlingEvent)

	if err = eh.eventService.CargoWasHandled(cargo); err != nil {
		return err
	}

	return eh.cargos.Store(cargo)
}
