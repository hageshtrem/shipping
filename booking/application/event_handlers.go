package application

import (
	"booking/domain"
	handling "booking/pb/handling/pb"

	"google.golang.org/protobuf/proto"
)

type EventHandler interface {
	Handle(event proto.Message) error
}

type cargoHandledEventHandler struct {
	cargos       domain.CargoRepository
	eventService EventService
}

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

func decodeHandlingActivity(activity *handling.Activity) domain.HandlingActivity {
	return domain.HandlingActivity{
		Type:         domain.HandlingEventType(activity.Type),
		Location:     domain.UNLocode(activity.Location),
		VoyageNumber: domain.VoyageNumber(activity.VoyageNumber),
	}
}
