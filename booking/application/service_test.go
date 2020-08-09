package application

import (
	"os"
	"testing"
	"time"

	"booking/domain"
	"booking/infrastructure"
)

var S Service
var routingServer string = "localhost:50051"

func init() {
	e := os.Getenv("ADDR")
	if e != "" {
		routingServer = e
	}
	rs, err := infrastructure.NewRoutingService(routingServer)
	if err != nil {
		panic(err)
	}
	S = NewService(infrastructure.NewCargoRepository(), infrastructure.NewLocationRepository(), rs)
}

func bookDefaultCargo() (domain.TrackingID, error) {
	return S.BookNewCargo(domain.SESTO, domain.USNYC, time.Now().AddDate(0, 1, 0))
}

func TestBookNewCargo(t *testing.T) {
	cargoID, err := bookDefaultCargo()
	checkErr(err, t)
	t.Logf("tracking ID: %s", cargoID)
}

func TestRequestPossibleRoutesForCargo(t *testing.T) {
	cargoID, err := bookDefaultCargo()
	checkErr(err, t)

	routes := S.RequestPossibleRoutesForCargo(cargoID)
	if len(routes) == 0 {
		t.Fatal("Empty routes")
	}

	t.Logf("Routes: %v", routes)
}

func TestAssignCargoToRoute(t *testing.T) {
	cargoID, err := bookDefaultCargo()
	checkErr(err, t)

	itinerary := domain.Itinerary{}
	err = S.AssignCargoToRoute(cargoID, itinerary)
	t.Logf("Assigning empty route: %v", err)

	routes := S.RequestPossibleRoutesForCargo(cargoID)
	err = S.AssignCargoToRoute(cargoID, routes[0])
	checkErr(err, t)
}

func TestChangeDestination(t *testing.T) {
	cargoID, err := bookDefaultCargo()
	checkErr(err, t)

	err = S.ChangeDestination(cargoID, "")
	t.Logf("Invalid destination error: %v", err)

	err = S.ChangeDestination(cargoID, domain.FIHEL)
	checkErr(err, t)
}

func TestLocations(t *testing.T) {
	locations := S.Locations()
	if len(locations) == 0 {
		t.Fatal("Empty locations")
	}
	t.Log(locations)
}

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
