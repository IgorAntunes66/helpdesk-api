package main

import (
	"helpdesk/pkg"
	"log"
	"net"
	"net/http"
	"os"

	pb "helpdesk/pkg/pb"
	"helpdesk/tickets-service/internal/handler"
	"helpdesk/tickets-service/internal/repository"

	ticketsGrpc "helpdesk/tickets-service/internal/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver do postgres para o migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

func main() {
	runMigrations()

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
		pb.RegisterTicketServiceServer(grpcServer, ticketServer)

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
	r.Post("/tickets/{id}/comments", apiServer.CreateCommentHandler)
	r.Delete("/tickets/{id}", apiServer.DeleteTicketHandler)

	log.Println("Servidor HTTP iniciado na porta 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", err)
	}

}

func runMigrations() {
	migrationDir := "file://db/migrations"
	dbURL := os.Getenv("CHAVEDB")
	m, err := migrate.New(migrationDir, dbURL)
	if err != nil {
		log.Fatalf("Erro ao criar a instancia de migração: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro ao cplicar migrações: %v", err)
	}

	log.Println("Migrações aplicadas com sucesso!")
}
