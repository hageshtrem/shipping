package application

import (
	pb "booking/pb"
)

type Server struct {
	pb.UnimplementedBookingServiceServer
}

func NewServer() Server {
	return Server{pb.UnimplementedBookingServiceServer{}}
}
