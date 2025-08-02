package main

import (
	"helpdesk/pkg"
	"log"
	"net"
	"net/http"

	pb "helpdesk/pkg/pb"
	"helpdesk/tickets-service/internal/handler"
	"helpdesk/tickets-service/internal/repository"

	ticketsGrpc "helpdesk/tickets-service/internal/grpc"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	db, err := pkg.ConectaDB()
	if err != nil {
		log.Fatalf("Erro ao iniciar o banco de dados: %v", err)
	}

	repo := repository.NewRepository(db)

	go func() {
		// 1 - Criar listener na porta 8081 para o gRPC
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalf("Falha ao escutar a porta gRPC: %v", err)
		}
		// 2 - Criar uma nova instancia do servidor gRPC
		grpcServer := grpc.NewServer()

		//3 - Instanciar a sua implementação do serviço de Ticket
		// (passando o repositorio como dependecia)
		ticketServer := ticketsGrpc.NewServer(pb.UnimplementedTicketServiceServer{}, repo)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Falha ao iniciar o servidor gRPC: %v", err)
		}
	}()


	apiServer := handler.NewApiServer(repo)

	r := chi.NewRouter()
	r.Post("/tickets", apiServer.CreateTicketHandler)
	r.Get("/health", handler.HealthCheckHandler)
	r.Get("/tickets", apiServer.ListTicketsHandler)
	r.Get("/tickets/{id}", apiServer.GetTicketHandler)
	r.Put("/tickets/{id}", apiServer.UpdateTicketHandler)
	r.Delete("/tickets/{id}", apiServer.DeleteTicketHandler)

	log.Println("Servidor HTTP iniciado na porta 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", err)
	}

}
