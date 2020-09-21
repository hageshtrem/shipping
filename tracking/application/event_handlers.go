package application

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type EventHandler interface {
	Handle(event proto.Message)
}

type newCargoBookedEventHandler struct {
	cargos CargoViewModelRepository
}

func (eh *newCargoBookedEventHandler) Handle(event proto.Message) {
	// process event
	fmt.Printf("%v\n", event)
}

func NewCargoBookedEventHandler(cargos CargoViewModelRepository) EventHandler {
	return &newCargoBookedEventHandler{cargos}
}
