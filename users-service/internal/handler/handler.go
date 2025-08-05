package handler

import (
	"context"
	"encoding/json"
	"errors"
	"helpdesk/pkg/pb"
	"helpdesk/users-service/internal/model"
	"helpdesk/users-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiServer struct {
	rep *repository.Repository
}

func NewApiServer(rep *repository.Repository) *ApiServer {
	return &ApiServer{
		rep: rep,
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func (api *ApiServer) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var usuario model.User
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}
	newID, err := api.rep.CreateUser(usuario)
	usuario.ID = newID
	if err != nil {
		http.Error(w, "Erro ao inserir o usuario no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(usuario)
	if err != nil {
		http.Error(w, "Erro ao codificar o usuario em JSON", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	usuarios, err := api.rep.FindAllUsers()
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(usuarios)
	if err != nil {
		http.Error(w, "Erro ao codificar a lista em JSON", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
		return
	}

	user, err := api.rep.FindUserByID(idInt)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao consultar o usuario no banco de dados", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Erro ao codificar o usuario em json", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
		return
	}

	var u model.User

	err = json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusInternalServerError)
		return
	}

	err = api.rep.UpdateUser(idInt, u)
	if err != nil {
		http.Error(w, "Erro ao atualizar o usuario no banco de dados", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
		return
	}

	err = api.rep.DeleteUser(idInt)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao excluir o usuario do banco de dados", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *ApiServer) CreateUserTicketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Erro ao conectar ao servidor gRPC", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewTicketServiceClient(conn)

	var ticket pb.CreateTicketRequest
	err = json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}

	client.CreateTicket(context.Background(), &ticket)
	w.WriteHeader(http.StatusOK)
}
