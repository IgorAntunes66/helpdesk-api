package main

import (
	"helpdesk/pkg"
	"helpdesk/tickets-service/internal/handler"
	"helpdesk/tickets-service/internal/repository"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := pkg.ConectaDB()
	if err != nil {
		log.Fatalf("Erro ao iniciar o banco de dados: %v", err)
	}

	repo := repository.NewRepository(db)
	apiServer := handler.NewApiServer(repo)

	r := chi.NewRouter()
	r.Get("/health", handler.HealthCheckHandler)
	r.Post("/tickets", apiServer.CreateTicketHandler)

	http.ListenAndServe(":8080", r)
}
