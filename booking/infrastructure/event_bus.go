package infrastructure

import (
	"booking/application"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type eventBus struct {
	*amqp.Connection
	*amqp.Channel
}

func NewEventBus(uri string) (application.EventBus, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := channel.ExchangeDeclare(
		"shipping", // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
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
		"shipping", // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			// Headers:   amqp.Table{},
			ContentType:  "application/x-protobuf; proto=" + routingKey, // TODO: change when standardized
			Body:         messageContent,
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func (eb *eventBus) Close() {
	eb.Channel.Close()
	eb.Connection.Close()
}
