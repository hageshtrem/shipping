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

	"pathfinder"
	"pathfinder/pb"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultPort   = ":50051"
	defaultLogDir = "/var/log/pathfinder"
)

func initLogger(dir string) {
	log.SetFormatter(&log.JSONFormatter{})

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModeDir|0666); err != nil {
			log.Fatalf("Failed to initialize log file %s", err)
		}
	}

	filename := filepath.Join(dir, fmt.Sprintf("pathfinder-%s.log", time.Now().Format("01.02.2006")))
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to initialize log file %s", err)
	}

	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)
}

func main() {
	port := envString("PORT", defaultPort)
	logDir := envString("LOG_DIR", defaultLogDir)

	initLogger(logDir)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error while binding port: %s", err)
	}
	defer lis.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	s := grpc.NewServer()
	mySrv := pathfinder.NewGRPCServer()
	mySrv = pathfinder.NewLoggingServer(log.WithField("component", "pathfinder_service"), mySrv)
	pb.RegisterPathfinderServiceServer(s, mySrv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error while starting server: %s", err)
		}
	}()
	log.Infof("listening at %s", port)
	<-stop

	log.Info("shutting down")
	s.GracefulStop()
	log.Info("terminated")
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
