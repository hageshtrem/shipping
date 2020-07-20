package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	app "booking/application"
	"booking/pb"
)

const (
	port = ":5051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := app.NewServer()

	gs := grpc.NewServer()

	pb.RegisterBookingServiceServer(gs, &s)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
