package main

import (
	"helpdesk/pkg"
	"helpdesk/users-service/internal/handler"
	"helpdesk/users-service/internal/repository"
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
	r.Get("/users", apiServer.ListUsersHandler)
	r.Get("/users/{id}", apiServer.GetUserHandler)
	r.Post("/users", apiServer.CreateUserHandler)
	r.Put("/users/{id}", apiServer.UpdateUserHandler)
	r.Delete("/users/{id}", apiServer.DeleteUserHandler)
	http.ListenAndServe(":8080", r)
}
