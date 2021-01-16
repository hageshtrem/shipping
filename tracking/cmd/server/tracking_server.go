package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	cargoWasHandledEH := application.NewCargoWasHandledEventHandler(cargos)

	errChan := make(chan error)

	eventBus, err := infrastructure.NewEventBus(envCfg.RabbitURI)
	checkErr(err)
	defer eventBus.Close()
	eventBus.NotifyError(errChan)

	checkErr(eventBus.Subscribe(&booking.NewCargoBooked{}, newCargoEH))
	checkErr(eventBus.Subscribe(&booking.CargoDestinationChanged{}, destChangedEH))
	checkErr(eventBus.Subscribe(&booking.CargoToRouteAssigned{}, routeAssignedEH))
	checkErr(eventBus.Subscribe(&booking.CargoWasHandled{}, cargoWasHandledEH))

	trackingSvc := application.NewService(cargos)

	lis, err := net.Listen("tcp", envCfg.Port)
	checkErr(err)

	s := application.NewServer(trackingSvc)
	gs := grpc.NewServer()
	defer gs.GracefulStop()

	pb.RegisterTrackingServiceServer(gs, s)
	log.Printf("Service started at %s", envCfg.Port)
	go func() {
		if err := gs.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Printf("Received %v signal, shutting down...", s)
	case err := <-errChan:
		log.Fatal(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
