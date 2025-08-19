package main

import (
	pkg "helpdesk/db"
	"helpdesk/tickets-service/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"helpdesk/tickets-service/internal/handler"
	"helpdesk/tickets-service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver do postgres para o migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	runMigrations()

	db, err := pkg.ConectaDB()
	if err != nil {
		log.Fatalf("Erro ao iniciar o banco de dados: %v", err)
	}

	var jobs = make(chan int64)

	for i := 0; i < 3; i++ {
		go RunWorkers(jobs)
	}

	repo := repository.NewRepository(db)

	apiServer := handler.NewApiServer(repo, jobs)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/tickets", apiServer.CreateTicketHandler)
		r.Post("/tickets/{id}/comments", apiServer.CreateCommentHandler)
		r.Get("/tickets/my-tickets", apiServer.GetMyTicketsHandler)
		r.Get("/tickets", apiServer.ListTicketsHandler)
		r.Get("/tickets/{id}", apiServer.GetTicketHandler)
		r.Get("/tickets/{id}/comments", apiServer.ListCommentsByTicketHandler)
		r.Get("/tickets/comments/users/{id}", apiServer.ListCommentsByUserHandler)
		r.Put("/tickets/{id}", apiServer.UpdateTicketHandler)
		r.Put("/tickets/comments/{id}", apiServer.UpdateCommentHandler)
		r.Patch("/tickets/{id}/status", apiServer.UpdateTicketStatusHandler)
		r.Delete("/tickets/{id}", apiServer.DeleteTicketHandler)
		r.Delete("/tickets/comments/{id}", apiServer.DeleteCommentHandler)
	})
	r.Get("/health", handler.HealthCheckHandler)

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

func RunWorkers(jobs <-chan int64) {
	for id := range jobs {
		log.Println("Simulando o envio de notificação para o ticket: ", id)
		time.Sleep(2 * time.Second)
	}
}
