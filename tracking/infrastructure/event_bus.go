package infrastructure

import (
	"fmt"
	"reflect"
	"strings"
	"tracking/application"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type EventBus interface {
	Subscribe(event proto.Message, eventHandler application.EventHandler) error
}

type eventBus struct {
	*amqp.Connection
	*amqp.Channel
	queueName string
}

func NewEventBus(uri string) (EventBus, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := channel.ExchangeDeclare(
		"booking", // name TODO: change to shipping
		"fanout",  // type TODO: change to direct
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &eventBus{conn, channel, q.Name}, nil
}

func (eb *eventBus) Subscribe(event proto.Message, eventHandler application.EventHandler) error {
	routingKey := strings.Split(reflect.TypeOf(&event).String(), ".")[1]

	if err := eb.QueueBind(
		eb.queueName, // queue name
		routingKey,   // routing key
		"booking",    // exchange // TODO: change to shipping
		false,
		nil); err != nil {
		return err
	}

	msgs, err := eb.Consume(
		eb.queueName, // queue
		"",           // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			message := reflect.ValueOf(event).Interface()

			if err := proto.Unmarshal(d.Body, message.(proto.Message)); err != nil {
				// log
				fmt.Println(err)
			} else {
				eventHandler.Handle(message.(proto.Message))
			}
		}
	}()

	return nil
}