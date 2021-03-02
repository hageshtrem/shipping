package application

import (
	"booking/domain"
	"time"

	log "github.com/sirupsen/logrus"
)

type loggingService struct {
	logger *log.Entry
	next   Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger *log.Entry, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) BookNewCargo(origin domain.UNLocode, destination domain.UNLocode, deadline time.Time) (id domain.TrackingID, err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method":           "BookNewCargo",
			"origin":           origin,
			"destination":      destination,
			"arrival_deadline": deadline,
			"took":             time.Since(begin),
			"err":              err,
		}).Info()
	}(time.Now())
	return s.next.BookNewCargo(origin, destination, deadline)
}

func (s *loggingService) RequestPossibleRoutesForCargo(id domain.TrackingID) []domain.Itinerary {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method":      "RequestPossibleRoutesForCargo",
			"tracking_id": id,
			"took":        time.Since(begin),
		}).Info()
	}(time.Now())
	return s.next.RequestPossibleRoutesForCargo(id)
}

func (s *loggingService) AssignCargoToRoute(id domain.TrackingID, itinerary domain.Itinerary) (err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method":      "AssignCargoToRoute",
			"tracking_id": id,
			"took":        time.Since(begin),
			"err":         err,
		}).Info()
	}(time.Now())
	return s.next.AssignCargoToRoute(id, itinerary)
}

func (s *loggingService) ChangeDestination(id domain.TrackingID, l domain.UNLocode) (err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method":      "ChangeDestination",
			"tracking_id": id,
			"destination": l,
			"took":        time.Since(begin),
			"err":         err,
		}).Info()
	}(time.Now())
	return s.next.ChangeDestination(id, l)
}

func (s *loggingService) LoadCargo(id domain.TrackingID) (c Cargo, err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method":      "LoadCargo",
			"tracking_id": id,
			"took":        time.Since(begin),
			"err":         err,
		}).Info()
	}(time.Now())
	return s.next.LoadCargo(id)
}

func (s *loggingService) Locations() []Location {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method": "Locations",
			"took":   time.Since(begin),
		}).Info()
	}(time.Now())
	return s.next.Locations()
}

func (s *loggingService) Cargos() []Cargo {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method": "Cargos",
			"took":   time.Since(begin),
		}).Info()
	}(time.Now())
	return s.next.Cargos()
}
