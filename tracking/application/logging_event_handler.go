package application

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type loggingEventHandler struct {
	logger *log.Entry
	next   EventHandler
}

// NewLoggingEventHandler returns a new instance of a logging EventHandler.
func NewLoggingEventHandler(logger *log.Entry, next EventHandler) EventHandler {
	return &loggingEventHandler{logger, next}
}

func (eh *loggingEventHandler) Handle(event proto.Message) (err error) {
	defer func(begin time.Time) {
		eh.logger.WithFields(log.Fields{
			"event": fmt.Sprintf("%v", event),
			"took":  time.Since(begin),
			"err":   err,
		}).Info()
	}(time.Now())
	return eh.next.Handle(event)
}
