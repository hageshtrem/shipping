package infrastructure

import (
	"booking/domain"
	"errors"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var address string = "localhost:50051"

func init() {
	e := os.Getenv("ADDR")
	if e != "" {
		address = e
	}
}

func TestRoutingService(t *testing.T) {
	s, err := NewRoutingService(address, log.NewEntry(log.New()))
	if err != nil {
		t.Fatal(err)
	}

	spec := domain.RouteSpecification{
		Origin:          domain.SESTO,
		Destination:     domain.AUMEL,
		ArrivalDeadline: time.Now().AddDate(0, 1, 0),
	}

	itinerary := s.FetchRoutesForSpecification(spec)
	if len(itinerary) == 0 {
		t.Fatal(errors.New("empty itinerary"))
	}
}
