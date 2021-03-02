package infrastructure

import (
	"booking/domain"
	"booking/pb/pathfinder/pb"
	pathfinderPb "booking/pb/pathfinder/pb"
	"context"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type proxyService struct {
	pathfinderPb.PathfinderServiceClient
	logger *log.Entry
}

// NewRoutingService returns implementation for domain.RoutingService.
func NewRoutingService(address string, logger *log.Entry) (domain.RoutingService, error) {
	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &proxyService{pb.NewPathfinderServiceClient(conn), logger}, nil
}

func (s *proxyService) FetchRoutesForSpecification(rs domain.RouteSpecification) []domain.Itinerary {
	req := pathfinderPb.ShortestPathReq{
		Origin:      string(rs.Origin),
		Destination: string(rs.Destination),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result := []domain.Itinerary{}
	stream, err := s.ShortestPath(ctx, &req)
	if err != nil {
		s.logger.Error(err)
		return result
	}

	for {
		path, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.logger.Error(err)
			return result
		}
		itinerary, err := assembly(path)
		if err != nil {
			s.logger.Error(err)
			return result
		}
		result = append(result, itinerary)
	}

	return result
}

func assembly(path *pathfinderPb.TransitPath) (domain.Itinerary, error) {
	legs := []domain.Leg{}
	for _, p := range path.GetEdges() {
		load, err := ptypes.Timestamp(p.GetDepartue())
		if err != nil {
			return domain.Itinerary{}, err
		}
		unload, err := ptypes.Timestamp(p.GetArrival())
		if err != nil {
			return domain.Itinerary{}, err
		}

		legs = append(legs, domain.Leg{
			VoyageNumber:   domain.VoyageNumber(p.GetVoyageNumber()),
			LoadLocation:   domain.UNLocode(p.GetOrigin()),
			UnloadLocation: domain.UNLocode(p.GetDestination()),
			LoadTime:       load,
			UnloadTime:     unload,
		})
	}
	return domain.Itinerary{Legs: legs}, nil
}
