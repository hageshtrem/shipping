package main

import (
	"context"
	"io"
	"log"
	"time"

	"pathfinder/pb"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	origin      = "SESTO"
	destination = "CNHKG"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPathfinderServiceClient(conn)

	// Contact the server and print out its response.
	req := pb.ShortestPathReq{
		Origin:      origin,
		Destination: destination,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.ShortestPath(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		path, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("%v\n\n", path)
	}
}
