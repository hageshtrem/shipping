package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codingconcepts/env"
	"google.golang.org/grpc"

	app "booking/application"
	infra "booking/infrastructure"
	"booking/pb"
)

type envConfig struct {
	Port        string `env:"PORT" default:":5051"`
	RoutingAddr string `env:"ROUTING_ADDR" default:"localhost:50051"`
	RabbitURI   string `env:"RABBIT_URI" default:"amqp://guest:guest@localhost:5672/"`
}

func main() {
	envCfg := envConfig{}
	checkErr(env.Set(&envCfg))

	routingSvc, err := infra.NewRoutingService(envCfg.RoutingAddr)
	checkErr(err)
	cargos := infra.NewCargoRepository()
	locations := infra.NewLocationRepository()
	eventBus, err := infra.NewEventBus(envCfg.RabbitURI)
	checkErr(err)
	defer eventBus.Close()
	eventService := app.NewEventService(eventBus)

	bookingSvc := app.NewService(cargos, locations, routingSvc, eventService)

	lis, err := net.Listen("tcp", envCfg.Port)
	checkErr(err)

	s := app.NewServer(bookingSvc)
	gs := grpc.NewServer()
	defer gs.GracefulStop()

	pb.RegisterBookingServiceServer(gs, s)
	errChan := make(chan error)
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
