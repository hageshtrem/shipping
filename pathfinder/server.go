package pathfinder

import (
	"pathfinder/pb"

	"github.com/golang/protobuf/ptypes"
	pf "github.com/marcusolsson/pathfinder"
	"github.com/marcusolsson/pathfinder/path"
)

type grpcServer struct {
	pb.UnimplementedPathfinderServiceServer
	pathService pf.PathService
}

// NewGRPCServer creates a new grpc server instance.
func NewGRPCServer() pb.PathfinderServiceServer {
	return &grpcServer{pathService: pf.NewPathService()}
}

func (gs *grpcServer) ShortestPath(req *pb.ShortestPathReq, stream pb.PathfinderService_ShortestPathServer) error {
	paths, err := gs.pathService.ShortestPath(req.Origin, req.Destination)
	if err != nil {
		return err
	}

	for _, p := range paths {
		pbPath, err := assembly(&p)
		if err != nil {
			return err
		}
		stream.Send(pbPath)
	}

	return nil
}

func assembly(in *path.TransitPath) (*pb.TransitPath, error) {
	edges := make([]*pb.TransitPath_TransitEdge, 0, len(in.Edges))

	for _, edge := range in.Edges {
		departure, err := ptypes.TimestampProto(edge.Departure)
		if err != nil {
			return nil, err
		}
		arrival, err := ptypes.TimestampProto(edge.Arrival)
		if err != nil {
			return nil, err
		}
		edges = append(edges, &pb.TransitPath_TransitEdge{
			VoyageNumber: edge.VoyageNumber,
			Origin:       edge.Origin,
			Destination:  edge.Destination,
			Departue:     departure,
			Arrival:      arrival,
		})
	}

	return &pb.TransitPath{Edges: edges}, nil
}
