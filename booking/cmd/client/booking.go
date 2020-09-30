package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/codingconcepts/env"
	cli "github.com/jawher/mow.cli"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"booking/pb"
)

const timeFormat = "02.01.2006"

type envConfig struct {
	Address string `env:"ADDR" default:"localhost:5051"`
}

func main() {
	conf := envConfig{}
	err := env.Set(&conf)
	checkErr(err)

	// Set up a connection to the server.
	conn, err := grpc.Dial(conf.Address, grpc.WithInsecure(), grpc.WithBlock())
	checkErr(err)
	defer conn.Close()

	c := pb.NewBookingServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app := cli.App("booking_client", "The client for booking service")
	app.Command("locations", "show available locations", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			locationResp, err := c.Locations(ctx, &empty.Empty{})
			checkErr(err)
			locations := locationResp.GetLocations()
			fmt.Println("Locations:")
			for _, loc := range locations {
				fmt.Printf("\t%s\n", loc)
			}
		}
	})

	app.Command("book", "book a new cargo", func(cmd *cli.Cmd) {
		var (
			origin      = cmd.StringArg("ORIG", "", "UNLOCODE of origin")
			destination = cmd.StringArg("DEST", "", "UNLOCODE of destination")
			deadline    = cmd.StringArg("DEADLINE", time.Now().AddDate(0, 1, 0).Format(timeFormat), "Deadline in format dd.mm.yyy")
		)
		cmd.Action = func() {
			t, err := time.Parse(timeFormat, *deadline)
			checkErr(err)
			resp, err := c.BookNewCargo(ctx, &pb.BookNewCargoRequest{
				Origin:      *origin,
				Destination: *destination,
				Deadline:    timestamppb.New(t),
			})
			checkErr(err)
			fmt.Printf("New cargo booked: %s\n", resp.GetTrackingId())
		}
	})

	app.Command("change_dest", "change destination of cargo", func(cmd *cli.Cmd) {
		var (
			trackingID  = cmd.StringArg("ID", "", "Tracking ID")
			destination = cmd.StringArg("DEST", "", "UNLOCODE of destination")
		)
		cmd.Action = func() {
			_, err := c.ChangeDestination(ctx, &pb.ChangeDestinationRequest{
				TrackingId:  *trackingID,
				Destination: *destination,
			})
			checkErr(err)
		}
		fmt.Println("Destination changed")
	})

	app.Command("list", "list cargos", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			cargoResp, err := c.Cargos(ctx, &empty.Empty{})
			checkErr(err)
			fmt.Println("Cargos:")
			for _, cargo := range cargoResp.GetCargos() {
				fmt.Printf("\tID: %s\tOrigin: %s\tDestination: %s\n", cargo.GetTrackingId(), cargo.GetOrigin(), cargo.GetDestination())
			}
		}
	})

	app.Command("show", "show info about specified cargo", func(cmd *cli.Cmd) {
		var trackingID = cmd.StringArg("ID", "", "Tracking ID")
		cmd.Action = func() {
			cargoResp, err := c.LoadCargo(ctx, &pb.LoadCargoRequest{
				TrackingId: *trackingID,
			})
			checkErr(err)
			cargo := cargoResp.GetCargo()
			fmt.Printf("ID: %s\nRouted: %t\nOrigin: %s\nDestination: %s\nMisrouted: %t\nArrival deadline: %s\n",
				cargo.GetTrackingId(),
				cargo.GetRouted(),
				cargo.GetOrigin(),
				cargo.GetDestination(),
				cargo.GetMisrouted(),
				cargo.GetArrivalDeadline().AsTime().Format(timeFormat),
			)
			fmt.Println("Itinerary:")
			printRoute(cargo.GetLegs())
		}
	})

	app.Command("route", "assign cargo to route", func(cmd *cli.Cmd) {
		var trackingID = cmd.StringArg("ID", "", "Tracking ID")
		cmd.Action = func() {
			stream, err := c.RequestPossibleRoutesForCargo(ctx, &pb.RequestPossibleRoutesForCargoRequest{
				TrackingId: *trackingID,
			})
			checkErr(err)

			fmt.Println("Waiting for routes...")
			routes := []*pb.Itinerary{}
			for {
				route, err := stream.Recv()
				if err == io.EOF {
					break
				}
				checkErr(err)
				routes = append(routes, route)
			}

			if len(routes) == 0 {
				fmt.Println("There are no available routes")
				return
			}

			fmt.Println("Select route")
			for i, route := range routes {
				fmt.Printf("[%d]\n", i)
				printRoute(route.GetLegs())
			}
			reader := bufio.NewReaderSize(os.Stdin, 1)
			choice, err := reader.ReadString('\n')
			checkErr(err)
			index, err := strconv.Atoi(choice[:len(choice)-1]) // exclude \n
			checkErr(err)
			if index >= len(routes) {
				fmt.Println("Wrong number")
				return
			}

			// new context
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err = c.AssignCargoToRoute(ctx, &pb.AssignCargoToRouteRequest{
				TrackingId: *trackingID,
				Itinerary:  routes[index],
			})
			checkErr(err)
			fmt.Println("Route assigned")
		}
	})

	checkErr(app.Run(os.Args))
}

func printRoute(itinerary []*pb.Leg) {
	for _, leg := range itinerary {
		fmt.Printf("\tVoyage number: %s\tLoad location: %s\tLoad time: %s\tUnload location: %s\tUnload time: %s\n",
			leg.GetVoyageNumber(),
			leg.GetLoadLocation(),
			leg.GetLoadTime().AsTime().Format(timeFormat),
			leg.GetUnloadLocation(),
			leg.GetUnloadTime().AsTime().Format(timeFormat),
		)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
