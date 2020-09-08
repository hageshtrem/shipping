package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	app "booking/application"
	infra "booking/infrastructure"
	"booking/pb"
)

const (
	PORT         = ":5051"
	ROUTING_ADDR = "localhost:50051"
	RABBIT_URI   = "amqp://guest:guest@localhost:5672/"
)

func main() {
	var (
		port        = envString("PORT", PORT)
		routingAddr = envString("ROUTING_ADDR", ROUTING_ADDR)
		rabbit_uri  = envString("RABBIT_URI", RABBIT_URI)
	)

	routingSvc, err := infra.NewRoutingService(routingAddr)
	checkErr(err)
	cargos := infra.NewCargoRepository()
	locations := infra.NewLocationRepository()
	eventBus, err := infra.NewEventBus(rabbit_uri)
	checkErr(err)

	bookingSvc := app.NewService(cargos, locations, routingSvc, eventBus)

	lis, err := net.Listen("tcp", port)
	checkErr(err)

	s := app.NewServer(bookingSvc)

	gs := grpc.NewServer()

	pb.RegisterBookingServiceServer(gs, s)
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
