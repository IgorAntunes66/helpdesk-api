package handler

import (
	"encoding/json"
	"errors"
	"helpdesk/users-service/internal/middleware"
	"helpdesk/users-service/internal/model"
	"helpdesk/users-service/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type ApiServer struct {
	rep model.UserRepository
}

// Dando problema ao chamar os metodos do BD
func NewApiServer(rep model.UserRepository) *ApiServer {
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
	w.WriteHeader(http.StatusOK)
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

	user, err := api.rep.FindUserByID(int64(idInt))
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao consultar o usuario no banco de dados", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

	idReq, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Não foi possivel extrair o ID do usuario do token", http.StatusInternalServerError)
		return
	}

	if int64(idInt) != idReq {
		http.Error(w, "Permissão não concedida", http.StatusForbidden)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusInternalServerError)
		return
	}

	err = api.rep.UpdateUser(int64(idInt), u)
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

	idReq, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Não foi possivel extrair o ID do usuario do token", http.StatusInternalServerError)
		return
	}

	if int64(idInt) != idReq {
		http.Error(w, "Permissão não concedida", http.StatusForbidden)
		return
	}

	err = api.rep.DeleteUser(int64(idInt))
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Usuario não encontrado no banco de dados", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao excluir o usuario do banco de dados", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *ApiServer) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	userDB, err := api.rep.FindUserByEmail(loginReq)
	if err != nil {
		http.Error(w, "Email ou senha", http.StatusNotFound)
		return
	}

	tokenJwt, err := utils.GerarToken(userDB)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			http.Error(w, "Token expirado", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Não foi possivel gerar o tokenJwt", http.StatusUnauthorized)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"token": tokenJwt}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Erro ao codigicar o token JWT", http.StatusInternalServerError)
		return
	}
}

func (s *ApiServer) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Não foi possivel extrair o ID do usuario do token", http.StatusInternalServerError)
		return
	}

	user, err := s.rep.FindUserByID(int64(userID))
	if err != nil {
		http.Error(w, "Erro ao encontrar o usuario no banco de dados", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Erro ao codificar o usuario", http.StatusInternalServerError)
		return
	}
}
