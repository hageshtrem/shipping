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

func encodeCargoToRouteAssigned(c *domain.Cargo) (*pb.CargoToRouteAssigned, error) {
	pbItinerary, err := encodeItinerary(&c.Itinerary)
	if err != nil {
		return nil, err
	}

	delivery, err := encodeDelivery(&c.Delivery)
	if err != nil {
		return nil, err
	}

	return &pb.CargoToRouteAssigned{
		TrackingId: string(c.TrackingID),
		Itinerary:  pbItinerary,
		Delivery:   delivery,
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
		Type:         domain.HandlingEventType(activity.GetType()),
		Location:     domain.UNLocode(activity.GetLocation()),
		VoyageNumber: domain.VoyageNumber(activity.GetVoyageNumber()),
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

func encodeItinerary(itinerary *domain.Itinerary) (*pb.Itinerary, error) {
	pbLegs, err := encodeLegs(itinerary.Legs)
	if err != nil {
		return nil, err
	}

	return &pb.Itinerary{Legs: pbLegs}, nil
}

func decodeItinerary(itinerary *pb.Itinerary) (*domain.Itinerary, error) {
	legs := make([]domain.Leg, 0, len(itinerary.Legs))

	for _, leg := range itinerary.Legs {
		loadTime, err := ptypes.Timestamp(leg.LoadTime)
		if err != nil {
			return nil, err
		}
		unloadTime, err := ptypes.Timestamp(leg.UnloadTime)
		if err != nil {
			return nil, err
		}
		legs = append(legs, domain.Leg{
			VoyageNumber:   domain.VoyageNumber(leg.VoyageNumber),
			LoadLocation:   domain.UNLocode(leg.LoadLocation),
			UnloadLocation: domain.UNLocode(leg.UnloadLocation),
			LoadTime:       loadTime,
			UnloadTime:     unloadTime,
		})
	}

	return &domain.Itinerary{Legs: legs}, nil
}

func encodeLegs(legs []domain.Leg) ([]*pb.Leg, error) {
	pbLegs := make([]*pb.Leg, 0, len(legs))

	for _, leg := range legs {
		loadTime, err := ptypes.TimestampProto(leg.LoadTime)
		if err != nil {
			return nil, err
		}
		unloadTime, err := ptypes.TimestampProto(leg.UnloadTime)
		if err != nil {
			return nil, err
		}
		pbLegs = append(pbLegs, &pb.Leg{
			VoyageNumber:   string(leg.VoyageNumber),
			LoadLocation:   string(leg.LoadLocation),
			UnloadLocation: string(leg.UnloadLocation),
			LoadTime:       loadTime,
			UnloadTime:     unloadTime,
		})
	}

	return pbLegs, nil
}

func encodeCargo(cargo *Cargo) (*pb.Cargo, error) {
	arrivalDeadline, err := ptypes.TimestampProto(cargo.ArrivalDeadline)
	if err != nil {
		return nil, err
	}
	pbLegs, err := encodeLegs(cargo.Legs)
	if err != nil {
		return nil, err
	}
	return &pb.Cargo{
		ArrivalDeadline: arrivalDeadline,
		Destination:     cargo.Destination,
		Legs:            pbLegs,
		Misrouted:       cargo.Misrouted,
		Origin:          cargo.Origin,
		Routed:          cargo.Routed,
		TrackingId:      cargo.TrackingID,
	}, nil
}
