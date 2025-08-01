package handler

import (
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
	"google.golang.org/protobuf/types/known/timestamppb"
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
	}
	newID, err := api.rep.CreateUser(usuario)
	usuario.ID = newID
	if err != nil {
		http.Error(w, "Erro ao inserir a tarefa no banco de dados", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(usuario)
	if err != nil {
		http.Error(w, "Erro ao codificar o usuario em JSON", http.StatusInternalServerError)
	}
}

func (api *ApiServer) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	usuarios, err := api.rep.FindAllUsers()
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(usuarios)
	if err != nil {
		http.Error(w, "Erro ao codificar a lista em JSON", http.StatusInternalServerError)
	}
}

func (api *ApiServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
	}

	user, err := api.rep.FindUserByID(idInt)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Erro ao consultar o usuario no banco de dados", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Erro ao codificar o usuario em json", http.StatusInternalServerError)
	}
}

func (api *ApiServer) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
	}

	var u model.User

	err = json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusInternalServerError)
	}

	err = api.rep.UpdateUser(idInt, u)
	if err != nil {
		http.Error(w, "Erro ao atualizar o usuario no banco de dados", http.StatusInternalServerError)
	}
}

func (api *ApiServer) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
	}

	err = api.rep.DeleteUser(idInt)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Erro ao excluir o usuario do banco de dados", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *ApiServer) TestGrpc(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Erro ao conectar ao servidor gRPC", http.StatusInternalServerError)
	}
	defer conn.Close()

	client := pb.NewTicketServiceClient(conn)

	client.CreateTicket(r.Context(), &pb.CreateTicketRequest{
		Titulo:          "Teste",
		Descricao:       "Teste",
		Status:          "Teste",
		Diagnostico:     "Teste",
		Solucao:         "Teste",
		Prioridade:      "Teste",
		DataAbertura:    timestamppb.Now(),
		DataFechamento:  timestamppb.Now(),
		DataAtualizacao: timestamppb.Now(),
	})
	w.WriteHeader(http.StatusOK)
}
