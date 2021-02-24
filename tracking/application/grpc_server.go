package application

import (
	"context"
	"tracking/pb"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	service Service
	pb.UnimplementedTrackingServiceServer
}

// NewServer creates a GRPC server.
func NewServer(s Service) pb.TrackingServiceServer {
	return &grpcServer{s, pb.UnimplementedTrackingServiceServer{}}
}

func (s *grpcServer) Track(_ context.Context, trackingID *pb.TrackingID) (*pb.Cargo, error) {
	c, err := s.service.Track(trackingID.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	pbCargo, err := encodeCargo(&c)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return pbCargo, nil
}

func encodeCargo(c *Cargo) (*pb.Cargo, error) {
	eta, err := ptypes.TimestampProto(c.ETA)
	if err != nil {
		return nil, err
	}

	arrivalDeadline, err := ptypes.TimestampProto(c.ArrivalDeadline)
	if err != nil {
		return nil, err
	}

	pbEvents := make([]*pb.Event, 0, len(c.Events))
	for _, event := range c.Events {
		pbEvents = append(pbEvents, &pb.Event{
			Description: event.Description,
			Expected:    event.Expected,
		})
	}

	return &pb.Cargo{
		TrackingId:           c.TrackingID,
		StatusText:           c.StatusText,
		Origin:               c.Origin,
		Destination:          c.Destination,
		Eta:                  eta,
		NextExpectedActivity: c.NextExpectedActivity,
		ArrivalDeadline:      arrivalDeadline,
		IsMisdirected:        c.IsMisdirected,
		Events:               pbEvents,
	}, nil
}
