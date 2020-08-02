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
	S = NewService(infrastructure.NewCargoRepository(), rs)
}

func TestBookNewCargo(t *testing.T) {
	cargoID, err := S.BookNewCargo(domain.SESTO, domain.USNYC, time.Now().AddDate(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tracking ID: %s", cargoID)
}

func TestRequestPossibleRoutesForCargo(t *testing.T) {
	cargoID, err := S.BookNewCargo(domain.SESTO, domain.USNYC, time.Now().AddDate(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	routes := S.RequestPossibleRoutesForCargo(cargoID)
	if len(routes) == 0 {
		t.Fatal("Empty routes")
	}
	t.Logf("Routes: %v", routes)
}
