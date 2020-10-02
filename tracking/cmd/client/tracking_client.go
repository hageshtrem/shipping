package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"tracking/pb"

	"google.golang.org/grpc"
)

const timeFormat = "02.01.2006"
const address = "localhost:5052"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	checkErr(err)
	defer conn.Close()
	client := pb.NewTrackingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(os.Args) <= 1 {
		fmt.Printf("Usage:\n\ttracking_client ID\n")
		os.Exit(1)
	}

	cargo, err := client.Track(ctx, &pb.TrackingID{Id: os.Args[1]})
	checkErr(err)

	fmt.Printf("TrackingID: %s\nStatus: %s\nOrigin: %s\nDestination: %s\nETA: %s\nNext expected activity: %s\nArrival deadline: %s\n",
		cargo.GetTrackingId(),
		cargo.GetStatusText(),
		cargo.GetOrigin(),
		cargo.GetDestination(),
		cargo.GetEta().AsTime().Format(timeFormat),
		cargo.GetNextExpectedActivity(),
		cargo.GetArrivalDeadline().AsTime().Format(timeFormat),
	)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
