package grpc

import (
	"context"
	pb "helpdesk/pkg/pb"
	"helpdesk/tickets-service/internal/model"
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
	ticket := model.Ticket{}
	ticket.Titulo = req.GetTitulo()
	ticket.Descricao = req.GetDescricao()
	ticket.Status = req.GetStatus()
	ticket.Diagnostico = req.GetDiagnostico()
	ticket.Solucao = req.GetSolucao()
	ticket.Prioridade = req.GetPrioridade()

	if req.GetDataAbertura() != nil && req.GetDataAbertura().IsValid() {
		ticket.DataAbertura = req.GetDataAbertura().AsTime()
	}

	if req.GetDataFechamento() != nil && req.GetDataFechamento().IsValid() {
		ticket.DataFechamento = req.GetDataFechamento().AsTime()
	}

	if req.DataAtualizacao != nil && req.DataAtualizacao.IsValid() {
		ticket.DataAtualizacao = req.DataAtualizacao.AsTime()
	}

	ticket.Anexos = req.Anexos
	ticket.Tags = req.Tags

	for _, historicoPb := range req.GetHistorico() {
		historico := model.Comentario{
			ID:        historicoPb.GetId(),
			Descricao: historicoPb.GetDescricao(),
			UserID:    historicoPb.GetUserId(),
			TicketID:  historicoPb.GetTicketId(),
		}

		if historicoPb.GetData() != nil && historicoPb.GetData().IsValid() {
			historico.Data = historicoPb.GetData().AsTime()
		}

		ticket.Historico = append(ticket.Historico, historico)
	}
	ticket.CategoriaID = req.GetCategoriaId()
	ticket.ResponsavelID = req.GetResponsavelId()
	ticket.UserID = req.GetUserId()

	id, err := s.rep.CreateTicket(ticket)
	if err != nil {
		return nil, err
	}

	return &pb.TicketResponse{Id: int64(id)}, nil
}
