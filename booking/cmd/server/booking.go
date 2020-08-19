package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	app "booking/application"
	infra "booking/infrastructure"
	"booking/pb"
)

const (
	port        = ":5051"
	routingAddr = "localhost:50051"
)

func main() {
	var (
		port        = envString("PORT", port)
		routingAddr = envString("ROUTING_ADDR", routingAddr)
	)

	routingSvc, err := infra.NewRoutingService(routingAddr)
	checkErr(err)
	cargos := infra.NewCargoRepository()
	locations := infra.NewLocationRepository()

	bookingSvc := app.NewService(cargos, locations, routingSvc)

	lis, err := net.Listen("tcp", port)
	checkErr(err)

	s := app.NewServer(bookingSvc)

	gs := grpc.NewServer()

	pb.RegisterBookingServiceServer(gs, s)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
