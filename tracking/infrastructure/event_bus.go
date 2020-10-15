package infrastructure

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"tracking/application"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	exchangeName = "shipping"
	queueName    = "tracking.queue"
)

// EventBus provides the ability to subscribe to an event.
type EventBus interface {
	// Subscribe registers a handler for specific type of event.
	Subscribe(event proto.Message, eventHandler application.EventHandler) error
	// NotifyError registers a channel to handle runtime errors.
	NotifyError(errChan chan<- error)
	// Close closes connection.
	Close()
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

type eventBus struct {
	*rabbitmq.Connection
	*rabbitmq.Channel
	consumer *consumer
}

// NewEventBus creates a new event bus.
func NewEventBus(uri string) (EventBus, error) {
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

func (eb *eventBus) NotifyError(errChan chan<- error) {
	eb.consumer.errChan = errChan
}

func (eb *eventBus) Close() {
	eb.consumer.enough <- struct{}{}
	eb.Channel.Close()
	eb.Connection.Close()
}
