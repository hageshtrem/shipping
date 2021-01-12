package application

import (
	"booking/domain"
	handling "booking/pb/handling/pb"

	"google.golang.org/protobuf/proto"
)

type EventHandler interface {
	Handle(event proto.Message) error
}

type cargoHandled struct {
	cargos domain.CargoRepository
}

func CargoHandled(cargos domain.CargoRepository) EventHandler {
	return &cargoHandled{cargos}
}

func (eh *cargoHandled) Handle(event proto.Message) error {
	handlingEvent := event.(*handling.HandlingEvent)
	domainHandlingEvent := domain.HandlingEvent{
		TrackingID: domain.TrackingID(handlingEvent.TrackingId),
		Activity:   decodeHandlingActivity(handlingEvent.Activity),
	}
	cargo, err := eh.cargos.Find(domain.TrackingID(handlingEvent.TrackingId))
	if err != nil {
		return err
	}
	cargo.DeriveDeliveryProgress(domainHandlingEvent)
	return nil
}

func decodeHandlingActivity(activity *handling.Activity) domain.HandlingActivity {
	return domain.HandlingActivity{
		Type:         domain.HandlingEventType(activity.Type),
		Location:     domain.UNLocode(activity.Location),
		VoyageNumber: domain.VoyageNumber(activity.VoyageNumber),
	}
}
