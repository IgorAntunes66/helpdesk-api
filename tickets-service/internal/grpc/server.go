package grpc

import (
	"context"
	pb "helpdesk/pkg/ticketpb"
	"helpdesk/tickets-service/internal/repository"
)

type Server struct {
	pb.UnimplementedTicketServiceServer
	rep *repository.Repository
}

func NewServer(pb pb.UnimplementedTicketServiceServer, rep *repository.Repository) *Server {
	return &Server{
		rep: rep,
	}
}

func (s *Server) CreateTicket(ctx context.Context, req *pb.CreateTicketRequest) (*pb.TicketResponse, error) {
	
}
