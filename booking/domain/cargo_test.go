package domain

import(
	"testing"
	"time"
	"fmt"
)

func ExampleNewCargo() {
	var id TrackingID = "testid"
	spec := RouteSpecification{
		Origin: SESTO,
		Destination: AUMEL,
		ArrivalDeadline: time.Now().AddDate(0, 1, 0),
	}

	c := NewCargo(id, spec)

	fmt.Println(c)
}

func TestNewCargo(t *testing.T) {
	var id TrackingID = "testid"
	spec := RouteSpecification{
		Origin: SESTO,
		Destination: AUMEL,
		ArrivalDeadline: time.Now().AddDate(0, 1, 0),
	}

	c := NewCargo(id, spec)

	t.Logf("Cargo: %v", c)
}