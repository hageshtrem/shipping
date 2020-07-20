package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"booking/pb"
)

const (
	address = "localhost:5051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookingServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BookNewCargo(ctx, &pb.BookNewCargoRequest{
		Origin:      "origin",
		Destination: "Destination",
		Deadline:    "deadline",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Booking: %s", r.String())
}
