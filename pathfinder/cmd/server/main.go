package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"pathfinder"
	"pathfinder/pb"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

const (
	defaultPort = ":50051"
)

func main() {
	port := envString("PORT", defaultPort)

	logger := log.NewJSONLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
	defer lis.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	s := grpc.NewServer()
	mySrv := pathfinder.NewGRPCServer()
	pb.RegisterPathfinderServiceServer(s, mySrv)
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}
	}()
	logger.Log("msg", "listening at "+port)
	<-stop

	logger.Log("msg", "shutting down")
	s.GracefulStop()
	logger.Log("msg", "terminated")
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
