package infrastructure

import (
	"booking/domain"
	"booking/pb/pathfinder/pb"
	pathfinderPb "booking/pb/pathfinder/pb"
	"context"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

type proxyService struct {
	pathfinderPb.PathfinderServiceClient
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
		log.Println(err)
		return result
	}

	for {
		path, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return result
		}
		itinerary, err := assembly(path)
		if err != nil {
			log.Println(err)
			return result
		}
		result = append(result, itinerary)
	}

	return result
}

// NewRoutingService returns implementation for domain.RoutingService.
func NewRoutingService(address string) (domain.RoutingService, error) {
	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &proxyService{pb.NewPathfinderServiceClient(conn)}, nil
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
