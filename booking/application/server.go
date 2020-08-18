package application

import (
	"booking/domain"
	pb "booking/pb"
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	service Service
	pb.UnimplementedBookingServiceServer
}

func NewServer(s Service) pb.BookingServiceServer {
	return &grpcServer{s, pb.UnimplementedBookingServiceServer{}}
}

func (s *grpcServer) BookNewCargo(_ context.Context, req *pb.BookNewCargoRequest) (*pb.BookNewCargoResponse, error) {
	origin := domain.UNLocode(req.GetOrigin())
	destination := domain.UNLocode(req.GetDestination())
	deadline, err := ptypes.Timestamp(req.GetDeadline())
	if err != nil {
		return nil, err
	}
	trackingID, err := s.service.BookNewCargo(origin, destination, deadline)
	if err != nil {
		return nil, err
	}
	return &pb.BookNewCargoResponse{TrackingId: string(trackingID)}, nil
}

func (s *grpcServer) RequestPossibleRoutesForCargo(req *pb.RequestPossibleRoutesForCargoRequest, stream pb.BookingService_RequestPossibleRoutesForCargoServer) error {
	id := domain.TrackingID(req.GetTrackingId())
	itineraries := s.service.RequestPossibleRoutesForCargo(id)

	for _, itin := range itineraries {
		pbItin, err := encodeItinerary(&itin)
		if err != nil {
			return err
		}
		stream.Send(pbItin)
	}

	return nil
}

func (s *grpcServer) AssignCargoToRoute(_ context.Context, req *pb.AssignCargoToRouteRequest) (*pb.Error, error) {
	id := domain.TrackingID(req.GetTrackingId())
	itinerary, err := decodeItinerary(req.GetItinerary())
	if err != nil {
		return nil, err
	}
	err = s.service.AssignCargoToRoute(id, *itinerary)
	if err != nil {
		return &pb.Error{Error: err.Error()}, nil
	}
	return nil, nil
}

func (s *grpcServer) ChangeDestination(_ context.Context, req *pb.ChangeDestinationRequest) (*pb.Error, error) {
	id := domain.TrackingID(req.GetTrackingId())
	dest := domain.UNLocode(req.GetDestination())
	if err := s.service.ChangeDestination(id, dest); err != nil {
		return &pb.Error{Error: err.Error()}, nil
	}
	return nil, nil
}

func (s *grpcServer) Locations(_ context.Context, _ *empty.Empty) (*pb.LocationsResponse, error) {
	locations := s.service.Locations()
	pbLocations := []*pb.LocationsResponse_Location{}
	for _, loc := range locations {
		pbLocations = append(pbLocations, &pb.LocationsResponse_Location{
			UnLocode: loc.UNLocode,
			Name:     loc.Name,
		})
	}
	return &pb.LocationsResponse{Locations: pbLocations}, nil
}

func (s *grpcServer) LoadCargo(_ context.Context, req *pb.LoadCargoRequest) (*pb.LoadCargoResponse, error) {
	id := domain.TrackingID(req.GetTrackingId())
	cargo, err := s.service.LoadCargo(id)
	if err != nil {
		return &pb.LoadCargoResponse{Result: &pb.LoadCargoResponse_Error{Error: err.Error()}}, err
	}
	pbCargo, err := encodeCargo(&cargo)
	if err != nil {
		return &pb.LoadCargoResponse{Result: &pb.LoadCargoResponse_Error{Error: err.Error()}}, err
	}

	return &pb.LoadCargoResponse{Result: &pb.LoadCargoResponse_Cargo{Cargo: pbCargo}}, nil
}

func (s *grpcServer) Cargos(_ context.Context, _ *empty.Empty) (*pb.CargosResponse, error) {
	cargos := s.service.Cargos()
	pbCargos := make([]*pb.Cargo, 0, len(cargos))
	for _, cargo := range cargos {
		pbCargo, err := encodeCargo(&cargo)
		if err != nil {
			return nil, err
		}
		pbCargos = append(pbCargos, pbCargo)
	}
	return &pb.CargosResponse{Cargos: pbCargos}, nil
}

func encodeItinerary(itinerary *domain.Itinerary) (*pb.Itinerary, error) {
	pbLegs, err := encodeLegs(itinerary.Legs)
	if err != nil {
		return nil, err
	}

	return &pb.Itinerary{Legs: pbLegs}, nil
}

func decodeItinerary(itinerary *pb.Itinerary) (*domain.Itinerary, error) {
	legs := make([]domain.Leg, 0, len(itinerary.Legs))

	for _, leg := range itinerary.Legs {
		loadTime, err := ptypes.Timestamp(leg.LoadTime)
		if err != nil {
			return nil, err
		}
		unloadTime, err := ptypes.Timestamp(leg.UnloadTime)
		if err != nil {
			return nil, err
		}
		legs = append(legs, domain.Leg{
			VoyageNumber:   domain.VoyageNumber(leg.VoyageNumber),
			LoadLocation:   domain.UNLocode(leg.LoadLocation),
			UnloadLocation: domain.UNLocode(leg.UnloadLocation),
			LoadTime:       loadTime,
			UnloadTime:     unloadTime,
		})
	}

	return &domain.Itinerary{Legs: legs}, nil
}

func encodeLegs(legs []domain.Leg) ([]*pb.Leg, error) {
	pbLegs := make([]*pb.Leg, 0, len(legs))

	for _, leg := range legs {
		loadTime, err := ptypes.TimestampProto(leg.LoadTime)
		if err != nil {
			return nil, err
		}
		unloadTime, err := ptypes.TimestampProto(leg.UnloadTime)
		if err != nil {
			return nil, err
		}
		pbLegs = append(pbLegs, &pb.Leg{
			VoyageNumber:   string(leg.VoyageNumber),
			LoadLocation:   string(leg.LoadLocation),
			UnloadLocation: string(leg.UnloadLocation),
			LoadTime:       loadTime,
			UnloadTime:     unloadTime,
		})
	}

	return pbLegs, nil
}

func encodeCargo(cargo *Cargo) (*pb.Cargo, error) {
	arrivalDeadline, err := ptypes.TimestampProto(cargo.ArrivalDeadline)
	if err != nil {
		return nil, err
	}
	pbLegs, err := encodeLegs(cargo.Legs)
	if err != nil {
		return nil, err
	}
	return &pb.Cargo{
		ArrivalDeadline: arrivalDeadline,
		Destination:     cargo.Destination,
		Legs:            pbLegs,
		Misrouted:       cargo.Misrouted,
		Origin:          cargo.Origin,
		Routed:          cargo.Routed,
		TrackingId:      cargo.TrackingID,
	}, nil
}
