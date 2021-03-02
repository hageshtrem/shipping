package application

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type loggingService struct {
	logger *log.Entry
	next   Service
}

// NewLoggingService returns a new instance of Service with logging decorator.
func NewLoggingService(logger *log.Entry, next Service) Service {
	return &loggingService{logger, next}
}

func (s *loggingService) Track(id string) (c Cargo, err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method": "Track",
			"id":     id,
			"took":   time.Since(begin),
			"err":    err,
		}).Info()
	}(time.Now())
	return s.next.Track(id)
}
