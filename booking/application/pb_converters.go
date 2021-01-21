package application

import (
	"booking/domain"
	"booking/pb"
	handling "booking/pb/handling/pb"

	"github.com/golang/protobuf/ptypes"
)

func encodeNewCargoBooked(c *domain.Cargo) (*pb.NewCargoBooked, error) {
	pbArrivalDeadline, err := ptypes.TimestampProto(c.RouteSpecification.ArrivalDeadline)
	if err != nil {
		return nil, err
	}
	delivery, err := encodeDelivery(&c.Delivery)
	if err != nil {
		return nil, err
	}
	return &pb.NewCargoBooked{
		TrackingId:      string(c.TrackingID),
		Origin:          string(c.Origin),
		Destination:     string(c.RouteSpecification.Destination),
		ArrivalDeadline: pbArrivalDeadline,
		Delivery:        delivery,
	}, nil
}

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
		TransportStatus:       encodeTransportStatus(d.TransportStatus),
		NextExpectedActivity:  encodeHandlingActivity(&d.NextExpectedActivity),
		LastEvent:             encodeHandlingEvent(&d.LastEvent),
		IsLastEventExpected:   d.Itinerary.IsExpected(d.LastEvent),
		LastKnownLocation:     string(d.LastKnownLocation),
		CurrentVoyage:         string(d.CurrentVoyage),
		Eta:                   eta,
		IsMisdirected:         d.IsMisdirected,
		IsUnloadAtDestination: d.IsUnloadedAtDestination,
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

func decodeHandlingActivity(activity *handling.Activity) domain.HandlingActivity {
	return domain.HandlingActivity{
		Type:         domain.HandlingEventType(activity.Type),
		Location:     domain.UNLocode(activity.Location),
		VoyageNumber: domain.VoyageNumber(activity.VoyageNumber),
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

func encodeTransportStatus(ts domain.TransportStatus) pb.TransportStatus {
	switch ts {
	case domain.NotReceived:
		return pb.TransportStatus_NotReceived
	case domain.InPort:
		return pb.TransportStatus_InPort
	case domain.OnboardCarrier:
		return pb.TransportStatus_OnboardCarrier
	case domain.Claimed:
		return pb.TransportStatus_Claimed
	case domain.Unknown:
		fallthrough
	default:
		return pb.TransportStatus_Unknown
	}
}
