package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/codingconcepts/env"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	app "booking/application"
	infra "booking/infrastructure"
	"booking/pb"
	handling "booking/pb/handling/pb"
)

type envConfig struct {
	Port        string `env:"PORT" default:":5051"`
	RoutingAddr string `env:"ROUTING_ADDR" default:"localhost:50051"`
	RabbitURI   string `env:"RABBIT_URI" default:"amqp://guest:guest@localhost:5672/"`
	LogDir      string `env:"LOG_DIR" default:"/var/log/booking"`
}

func initLogger(dir string) {
	log.SetFormatter(&log.JSONFormatter{})

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModeDir|0666); err != nil {
			log.Fatalf("Failed to initialize log file %s", err)
		}
	}

	filename := filepath.Join(dir, fmt.Sprintf("booking-%s.log", time.Now().Format("01.02.2006")))
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to initialize log file %s", err)
	}

	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)
}

func main() {
	envCfg := envConfig{}
	checkErr(env.Set(&envCfg))

	initLogger(envCfg.LogDir)

	routingSvc, err := infra.NewRoutingService(envCfg.RoutingAddr, log.WithField("component", "routing_service"))
	checkErr(err)
	cargos := infra.NewCargoRepository()
	locations := infra.NewLocationRepository()

	eventBus, err := infra.NewEventBus(envCfg.RabbitURI, log.WithField("component", "event_bus"))
	checkErr(err)
	defer eventBus.Close()
	eventService := app.NewEventService(eventBus)
	cargoHandledEH := app.NewCargoHandledEventHandler(cargos, eventService)
	cargoHandledEH = app.NewLoggingEventHandler(
		log.WithFields(log.Fields{
			"component": "booking_event_handler",
			"handler":   "CargoHandledEventHandler"}),
		cargoHandledEH,
	)
	checkErr(eventBus.Subscribe(&handling.HandlingEvent{}, cargoHandledEH))

	bookingSvc := app.NewService(cargos, locations, routingSvc, eventService)
	bookingSvc = app.NewLoggingService(log.WithField("component", "booking_service"), bookingSvc)

	lis, err := net.Listen("tcp", envCfg.Port)
	checkErr(err)

	s := app.NewServer(bookingSvc)
	gs := grpc.NewServer()
	defer gs.GracefulStop()

	pb.RegisterBookingServiceServer(gs, s)
	errChan := make(chan error)
	log.Infof("Service started at %s", envCfg.Port)
	go func() {
		if err := gs.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Infof("Received %v signal, shutting down...", s)
	case err := <-errChan:
		log.Fatal(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
