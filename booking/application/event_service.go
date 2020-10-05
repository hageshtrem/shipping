package application

import (
	"booking/domain"
	"booking/pb"

	"github.com/golang/protobuf/ptypes"
	// "google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

type EventBus interface {
	Publish(proto.Message) error
}

type EventService interface {
	NewCargoBooked(*domain.Cargo) error
	DestinationChanged(*domain.Cargo) error
	CargoToRouteAssigned(*domain.Cargo) error
}

func NewEventService(eventBus EventBus) EventService {
	return &eventService{eventBus}
}

type eventService struct {
	eventBus EventBus
}

func (es *eventService) NewCargoBooked(c *domain.Cargo) error {
	pbArrivalDeadline, err := ptypes.TimestampProto(c.RouteSpecification.ArrivalDeadline)
	if err != nil {
		return err
	}

	event := &pb.NewCargoBooked{
		TrackingId:      string(c.TrackingID),
		Origin:          string(c.Origin),
		Destination:     string(c.RouteSpecification.Destination),
		ArrivalDeadline: pbArrivalDeadline,
	}

	return es.eventBus.Publish(event)
}

func (es *eventService) DestinationChanged(c *domain.Cargo) error {
	event := &pb.CargoDestinationChanged{
		TrackingId:  string(c.TrackingID),
		Destination: string(c.RouteSpecification.Destination),
	}

	return es.eventBus.Publish(event)
}

func (es *eventService) CargoToRouteAssigned(c *domain.Cargo) error {
	pbItinerary, err := encodeItinerary(&c.Itinerary)
	if err != nil {
		return err
	}

	pbEta, err := ptypes.TimestampProto(c.Delivery.ETA)
	if err != nil {
		return err
	}

	event := &pb.CargoToRouteAssigned{
		TrackingId: string(c.TrackingID),
		Eta:        pbEta,
		Itinerary:  pbItinerary,
	}

	return es.eventBus.Publish(event)
}
