package pathfinder

import (
	"fmt"
	"pathfinder/pb"
	"time"

	log "github.com/sirupsen/logrus"
)

type loggingServer struct {
	logger *log.Entry
	pb.PathfinderServiceServer
}

// NewLoggingServer returns a PathfinderSrviceServer wrapped with a logging decorator.
func NewLoggingServer(logger *log.Entry, next pb.PathfinderServiceServer) pb.PathfinderServiceServer {
	return &loggingServer{logger, next}
}

func (s *loggingServer) ShortestPath(req *pb.ShortestPathReq, stream pb.PathfinderService_ShortestPathServer) (err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method": "ShortestPath",
			"req":    fmt.Sprintf("%v", req),
			"took":   time.Since(begin),
			"err":    err,
		}).Info()
	}(time.Now())
	return s.PathfinderServiceServer.ShortestPath(req, stream)
}
