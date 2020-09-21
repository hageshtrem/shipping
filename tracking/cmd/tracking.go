package main

import (
	"log"
	"os"
	"tracking/application"
	"tracking/infrastructure"
	booking "tracking/pb/booking/pb"
)

const (
	// PORT         = ":5051"
	// ROUTING_ADDR = "localhost:50051"
	RABBIT_URI = "amqp://guest:guest@localhost:5672/"
)

func main() {
	var (
		// port        = envString("PORT", PORT)
		// routingAddr = envString("ROUTING_ADDR", ROUTING_ADDR)
		rabbit_uri = envString("RABBIT_URI", RABBIT_URI)
	)

	cargos := infrastructure.NewCargoViewModelRepository()
	newCargoEH := application.NewCargoBookedEventHandler(cargos)

	eventBus, err := infrastructure.NewEventBus(rabbit_uri)
	checkErr(err)
	err = eventBus.Subscribe(&booking.NewCargoBooked{}, newCargoEH)
	checkErr(err)

	forever := make(chan bool)
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
