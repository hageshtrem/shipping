package infrastructure

import (
	"booking/application"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	exchangeName = "shipping"
	queueName    = "booking.queue"
)

type eventBus struct {
	*rabbitmq.Connection
	*rabbitmq.Channel
	consumer *consumer
}

type handlerFunc func(body []byte) error

type consumer struct {
	handlersSync sync.RWMutex
	handlers     map[string]handlerFunc
	errChan      chan<- error
	enough       chan struct{}
}

func (con *consumer) process(msgs <-chan amqp.Delivery) {
	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					break
				}

				con.handlersSync.Lock()
				h, ok := con.handlers[d.Type]
				con.handlersSync.Unlock()
				if !ok {
					// There is no handler registered for the received message yet.
					// So just return the message to the queue.
					log.Printf("No handler for received message: %s", d.Type)
					_ = d.Nack(false, true)
					continue
				}

				if err := h(d.Body); err != nil {
					_ = d.Nack(false, true)
					con.processErr(err)
				}

				if err := d.Ack(false); err != nil {
					con.processErr(err)
				}
			case <-con.enough:
				break
			}
		}
	}()
}

func (con *consumer) addHandler(msgType string, handler handlerFunc) {
	con.handlersSync.Lock()
	defer con.handlersSync.Unlock()
	con.handlers[msgType] = handler
}

func (con *consumer) processErr(err error) {
	if con.errChan == nil {
		panic(err)
	}

	con.errChan <- err
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
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	if err := channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	consumer := consumer{
		handlersSync: sync.RWMutex{},
		handlers:     make(map[string]handlerFunc),
		errChan:      nil,
		enough:       make(chan struct{}),
	}
	consumer.process(msgs)

	return &eventBus{conn, channel, &consumer}, nil
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

func (eb *eventBus) Subscribe(event proto.Message, eventHandler application.EventHandler) error {
	routingKey := strings.Split(reflect.TypeOf(event).String(), ".")[1]

	if err := eb.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil); err != nil {
		return err
	}

	eb.consumer.addHandler(routingKey, func(body []byte) error {
		message := reflect.ValueOf(event).Interface()

		unmarshalOptions := proto.UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		}
		if err := unmarshalOptions.Unmarshal(body, message.(proto.Message)); err != nil {
			return err
		}

		if err := eventHandler.Handle(message.(proto.Message)); err != nil {
			return err
		}

		return nil
	})

	return nil
}

// Close closes connection.
func (eb *eventBus) Close() {
	eb.Channel.Close()
	eb.Connection.Close()
}
