package main

import (
	"log"
	"net"
	"os"
	"tracking/application"
	"tracking/infrastructure"
	"tracking/pb"
	booking "tracking/pb/booking/pb"

	"github.com/codingconcepts/env"
	"google.golang.org/grpc"
)

type envConfig struct {
	Port      string `env:"PORT" default:":5052"`
	RabbitURI string `env:"RABBIT_URI" default:"amqp://guest:guest@localhost:5672/"`
}

func main() {
	envCfg := envConfig{}
	checkErr(env.Set(&envCfg))

	cargos := infrastructure.NewCargoViewModelRepository()
	newCargoEH := application.NewCargoBookedEventHandler(cargos)
	destChangedEH := application.NewCargoDestinationChangedEventHandler(cargos)
	routeAssignedEH := application.NewCargoToRouteAssignedEventHandler(cargos)

	eventBus, err := infrastructure.NewEventBus(envCfg.RabbitURI)
	checkErr(err)
	checkErr(eventBus.Subscribe(&booking.NewCargoBooked{}, newCargoEH))
	checkErr(eventBus.Subscribe(&booking.CargoDestinationChanged{}, destChangedEH))
	checkErr(eventBus.Subscribe(&booking.CargoToRouteAssigned{}, routeAssignedEH))

	trackingSvc := application.NewService(cargos)

	lis, err := net.Listen("tcp", envCfg.Port)
	checkErr(err)

	s := application.NewServer(trackingSvc)
	gs := grpc.NewServer()

	pb.RegisterTrackingServiceServer(gs, s)
	log.Printf("Service started at %s", envCfg.Port)
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
