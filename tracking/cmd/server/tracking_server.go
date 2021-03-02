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
	app "tracking/application"
	"tracking/infrastructure"
	"tracking/pb"
	booking "tracking/pb/booking/pb"

	"github.com/codingconcepts/env"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type envConfig struct {
	Port      string `env:"PORT" default:":5052"`
	RabbitURI string `env:"RABBIT_URI" default:"amqp://guest:guest@localhost:5672/"`
	LogDir    string `env:"LOG_DIR" default:"/var/log/tracking"`
}

func initLogger(dir string) {
	log.SetFormatter(&log.JSONFormatter{})

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModeDir|0666); err != nil {
			log.Fatalf("Failed to initialize log file %s", err)
		}
	}

	filename := filepath.Join(dir, fmt.Sprintf("tracking-%s.log", time.Now().Format("01.02.2006")))
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

	ehLogger := log.WithField("component", "tracking_event_handler")
	cargos := infrastructure.NewCargoViewModelRepository()
	cargoBookedEH := app.NewCargoBookedEventHandler(cargos)
	cargoBookedEH = app.NewLoggingEventHandler(ehLogger.WithField("handler", "CargoBookedEH"), cargoBookedEH)
	destChangedEH := app.NewCargoDestinationChangedEventHandler(cargos)
	destChangedEH = app.NewLoggingEventHandler(ehLogger.WithField("handler", "CargoDestinationChangedEH"), destChangedEH)
	routeAssignedEH := app.NewCargoToRouteAssignedEventHandler(cargos)
	routeAssignedEH = app.NewLoggingEventHandler(ehLogger.WithField("handler", "CargoToRouteAssignedEH"), routeAssignedEH)
	cargoWasHandledEH := app.NewCargoWasHandledEventHandler(cargos)
	cargoWasHandledEH = app.NewLoggingEventHandler(ehLogger.WithField("handler", "CargoWasHandledEH"), cargoWasHandledEH)

	errChan := make(chan error)

	eventBus, err := infrastructure.NewEventBus(envCfg.RabbitURI, log.WithField("component", "event_bus"))
	checkErr(err)
	defer eventBus.Close()
	eventBus.NotifyError(errChan)

	checkErr(eventBus.Subscribe(&booking.NewCargoBooked{}, cargoBookedEH))
	checkErr(eventBus.Subscribe(&booking.CargoDestinationChanged{}, destChangedEH))
	checkErr(eventBus.Subscribe(&booking.CargoToRouteAssigned{}, routeAssignedEH))
	checkErr(eventBus.Subscribe(&booking.CargoWasHandled{}, cargoWasHandledEH))

	trackingSvc := app.NewService(cargos)
	trackingSvc = app.NewLoggingService(log.WithField("component", "tracking_service"), trackingSvc)

	lis, err := net.Listen("tcp", envCfg.Port)
	checkErr(err)

	s := app.NewServer(trackingSvc)
	gs := grpc.NewServer()
	defer gs.GracefulStop()

	pb.RegisterTrackingServiceServer(gs, s)
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
