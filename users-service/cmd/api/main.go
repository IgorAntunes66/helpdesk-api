package main

import (
	"fmt"
	pkg "helpdesk/db"
	"helpdesk/users-service/internal/handler"
	"helpdesk/users-service/internal/repository"
	"helpdesk/users-service/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	runMigrations()

	db, err := pkg.ConectaDB()
	if err != nil {
		log.Fatalf("Erro ao iniciar o banco de dados: %v", err)
	}

	repo := repository.NewRepository(db)
	apiServer := handler.NewApiServer(repo)

	r := chi.NewRouter()
	r.Get("/health", handler.HealthCheckHandler)
	r.Post("/users", apiServer.CreateUserHandler)
	r.Post("/users/login", apiServer.LoginUserHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/users/me", apiServer.GetMeHandler)
		r.Get("/users", apiServer.ListUsersHandler)
		r.Get("/users/{id}", apiServer.GetUserHandler)
		r.Put("/users/{id}", apiServer.UpdateUserHandler)
		r.Delete("/users/{id}", apiServer.DeleteUserHandler)
	})
	http.ListenAndServe(":8082", r)
	fmt.Println("Servidor HTTP iniciado na porta 8082")
}

func runMigrations() {
	migrationDir := "file://db/migrations"
	dbURL := os.Getenv("CHAVEDB")
	m, err := migrate.New(migrationDir, dbURL)
	if err != nil {
		log.Fatalf("Erro ao criar a instancia de migração: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro ao aplicar migrações: %v", err)
	}

	log.Println("Migrações aplicadas com sucesso!")
}
