package application

import (
	"testing"
	"time"

	"booking/domain"
	"booking/infrastructure"
)

var S Service

func init() {
	S = NewService( infrastructure.NewCargoRepository() )
}

func TestBookNewCargo(t *testing.T) {
	cargoID, err := S.BookNewCargo(domain.SESTO, domain.USNYC, time.Now().AddDate(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tracking ID: %s", cargoID)
}
