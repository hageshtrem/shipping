package infrastructure

import (
	"booking/application"
	"fmt"
	"reflect"
	"strings"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	exchangeName = "shipping"
)

type eventBus struct {
	*rabbitmq.Connection
	*rabbitmq.Channel
}

// NewEventBus returns an implementation of application.EventBus.
func NewEventBus(uri string) (application.EventBus, error) {
	conn, err := rabbitmq.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := channel.ExchangeDeclare(
		exchangeName,        // name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	); err != nil {
		return nil, err
	}

	return &eventBus{conn, channel}, nil
}

func (eb *eventBus) Publish(message proto.Message) error {
	routingKey := strings.Split(reflect.TypeOf(message).String(), ".")[1]
	messageContent, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	if err := eb.Channel.Publish(
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/x-protobuf", // TODO: change when standardized
			Type:         routingKey,
			Body:         messageContent,
			DeliveryMode: amqp.Persistent,
			Priority:     0, // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// Close closes connection.
func (eb *eventBus) Close() {
	eb.Channel.Close()
	eb.Connection.Close()
}
