package application

import (
	"booking/domain"
	"booking/pb"

	"github.com/golang/protobuf/ptypes"
)

func encodeCargoWasHandled(c *domain.Cargo) (*pb.CargoWasHandled, error) {
	delivery, err := encodeDelivery(&c.Delivery)
	if err != nil {
		return nil, err
	}
	event := &pb.CargoWasHandled{
		TrackingId: string(c.TrackingID),
		Delivery:   delivery,
	}
	return event, nil
}

func encodeDelivery(d *domain.Delivery) (*pb.Delivery, error) {
	eta, err := ptypes.TimestampProto(d.ETA)
	if err != nil {
		return nil, err
	}

	pbDelivery := &pb.Delivery{
		TransportStatus:      d.TransportStatus.String(),
		NextExpectedActivity: encodeHandlingActivity(&d.NextExpectedActivity),
		Eta:                  eta,
		LastEvent:            encodeHandlingEvent(&d.LastEvent),
	}

	return pbDelivery, nil
}

func encodeHandlingEvent(e *domain.HandlingEvent) *pb.HandlingEvent {
	return &pb.HandlingEvent{
		TrackingId: string(e.TrackingID),
		Activity:   encodeHandlingActivity(&e.Activity),
	}
}

func encodeHandlingActivity(a *domain.HandlingActivity) *pb.HandlingActivity {
	return &pb.HandlingActivity{
		Type:         *encodeHandlingEventType(&a.Type),
		Location:     string(a.Location),
		VoyageNumber: string(a.VoyageNumber),
	}
}

func encodeHandlingEventType(et *domain.HandlingEventType) *pb.HandlingEventType {
	var result pb.HandlingEventType
	switch *et {
	case domain.NotHandled:
		result = pb.HandlingEventType_NotHandled
	case domain.Load:
		result = pb.HandlingEventType_Load
	case domain.Unload:
		result = pb.HandlingEventType_Unload
	case domain.Receive:
		result = pb.HandlingEventType_Receive
	case domain.Claim:
		result = pb.HandlingEventType_Claim
	case domain.Customs:
		result = pb.HandlingEventType_Customs
	}
	return &result
}
