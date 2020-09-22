package main

import (
	"log"
	"net"
	"os"
	"tracking/application"
	"tracking/infrastructure"
	"tracking/pb"
	booking "tracking/pb/booking/pb"

	"google.golang.org/grpc"
)

const (
	PORT       = ":5052"
	RABBIT_URI = "amqp://guest:guest@localhost:5672/"
)

func main() {
	var (
		port       = envString("PORT", PORT)
		rabbit_uri = envString("RABBIT_URI", RABBIT_URI)
	)

	cargos := infrastructure.NewCargoViewModelRepository()
	newCargoEH := application.NewCargoBookedEventHandler(cargos)

	eventBus, err := infrastructure.NewEventBus(rabbit_uri)
	checkErr(err)
	err = eventBus.Subscribe(&booking.NewCargoBooked{}, newCargoEH)
	checkErr(err)

	trackingSvc := application.NewService(cargos)

	lis, err := net.Listen("tcp", port)
	checkErr(err)

	s := application.NewServer(trackingSvc)
	gs := grpc.NewServer()

	pb.RegisterTrackingServiceServer(gs, s)
	log.Printf("Service started at %s", port)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
