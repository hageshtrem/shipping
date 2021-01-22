package application

import (
	"booking/domain"
	"booking/pb"

	"google.golang.org/protobuf/proto"
)

// EventBus is abstraction for Pub/Sub bus.
type EventBus interface {
	// Publish publishes the event to the bus.
	Publish(proto.Message) error
	// Subscribe registers a handler for specific type of event.
	Subscribe(event proto.Message, eventHandler EventHandler) error
	// Close closes the connection.
	Close()
}

// EventService provides methods for publish specific events.
type EventService interface {
	NewCargoBooked(*domain.Cargo) error
	DestinationChanged(*domain.Cargo) error
	CargoToRouteAssigned(*domain.Cargo) error
	CargoWasHandled(*domain.Cargo) error
}

// NewEventService returns an instance of EventService.
func NewEventService(eventBus EventBus) EventService {
	return &eventService{eventBus}
}

type eventService struct {
	eventBus EventBus
}

func (es *eventService) NewCargoBooked(c *domain.Cargo) error {
	event, err := encodeNewCargoBooked(c)
	if err != nil {
		return err
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
	event, err := encodeCargoToRouteAssigned(c)
	if err != nil {
		return err
	}

	return es.eventBus.Publish(event)
}

func (es *eventService) CargoWasHandled(c *domain.Cargo) error {
	event, err := encodeCargoWasHandled(c)
	if err != nil {
		return err
	}

	return es.eventBus.Publish(event)
}
