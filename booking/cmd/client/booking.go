package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"booking/pb"
)

const (
	address = "localhost:5051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	checkErr(err)
	defer conn.Close()
	c := pb.NewBookingServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	locationResp, err := c.Locations(ctx, &empty.Empty{})
	checkErr(err)
	locations := locationResp.GetLocations()
	log.Printf("Locations: %v", locations)
	log.Printf("Locations[0]: %v", locations[0])
	if len(locations) < 3 {
		log.Fatal("Not enought locations")
	}

	cargoID, err := c.BookNewCargo(ctx, &pb.BookNewCargoRequest{
		Origin:      locations[0].UnLocode,
		Destination: locations[1].UnLocode,
		Deadline:    timestamppb.Now(),
	})
	checkErr(err)
	log.Printf("BookNewCargo: %s", cargoID.GetTrackingId())

	pbErr, err := c.ChangeDestination(ctx, &pb.ChangeDestinationRequest{
		TrackingId:  cargoID.GetTrackingId(),
		Destination: locations[3].GetUnLocode(),
	})
	checkErr(err)
	log.Printf("ChangeDestination: %v", pbErr)

	stream, err := c.RequestPossibleRoutesForCargo(ctx, &pb.RequestPossibleRoutesForCargoRequest{
		TrackingId: cargoID.GetTrackingId(),
	})
	checkErr(err)
	routes := []*pb.Itinerary{}
	for {
		route, err := stream.Recv()
		if err == io.EOF {
			log.Println("Got EOF")
			break
		}
		checkErr(err)
		routes = append(routes, route)
		log.Printf("RequestPossibleRoutesForCargo: Route: %v", route)
	}

	if len(routes) != 0 {
		pbErr, err := c.AssignCargoToRoute(ctx, &pb.AssignCargoToRouteRequest{
			TrackingId: cargoID.GetTrackingId(),
			Itinerary:  routes[0],
		})
		checkErr(err)
		log.Printf("AssignCargoToRoute: %v", pbErr)
	}

	cargos, err := c.Cargos(ctx, &empty.Empty{})
	checkErr(err)
	log.Printf("Cargos: %v", cargos)

	cargoResp, err := c.LoadCargo(ctx, &pb.LoadCargoRequest{
		TrackingId: cargoID.GetTrackingId(),
	})
	checkErr(err)
	log.Printf("LoadCargo: %v", cargoResp)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
