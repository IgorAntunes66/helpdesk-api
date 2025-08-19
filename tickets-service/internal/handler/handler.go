package handler

import (
	"encoding/json"
	"fmt"
	"helpdesk/tickets-service/internal/model"
	"helpdesk/tickets-service/internal/repository"
	"helpdesk/tickets-service/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ApiServer struct {
	rep  *repository.Repository
	jobs chan int64
}

func NewApiServer(rep *repository.Repository, jobs chan int64) *ApiServer {
	return &ApiServer{
		rep:  rep,
		jobs: jobs,
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]any{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func (api *ApiServer) CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket model.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	userIdReq := r.Context().Value(middleware.UserIDKey).(int64)
	ticket.UserID = userIdReq

	id, err := api.rep.CreateTicket(ticket)
	if err != nil {
		http.Error(w, "Erro ao adicionar o ticket no banco de dados", http.StatusBadRequest)
		return
	}
	ticket.ID = id

	url := fmt.Sprintf("http://users-service:8082/users/%d", ticket.UserID)
	fmt.Printf("INFO: Serviço de tickets fazendo uma requisição interna para: %s\n", url)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Header de autorização ausente", http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		http.Error(w, "Header de autorização mal formatado", http.StatusUnauthorized)
		return
	}

	tokenString := headerParts[1]

	cliente := &http.Client{
		Timeout: 5 * time.Second,
	}

	reqInternal, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar requisição interna: "+err.Error(), http.StatusInternalServerError)
		return
	}

	reqInternal.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	reqInternal.Header.Set("Content-Type", "application/json")

	resposta, err := cliente.Do(reqInternal)
	if err != nil {
		http.Error(w, "ERRO: Falha ao fazer a requisição para o serviço de usuarios", http.StatusBadRequest)
		return
	}
	defer resposta.Body.Close()

	if resposta.StatusCode != http.StatusOK {
		http.Error(w, "Erro ao consultar serviço de usuarios", http.StatusBadRequest)
		return
	}

	var ticketAuthor model.TicketAuthor
	if err = json.NewDecoder(resposta.Body).Decode(&ticketAuthor); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da resposta do users-service", http.StatusBadRequest)
		return
	}

	ticket.Author = ticketAuthor

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ticket)
	if err != nil {
		http.Error(w, "Erro ao codificar o ticket em json", http.StatusInternalServerError)
		return
	}

	api.jobs <- ticket.ID
}

func (api *ApiServer) ListTicketsHandler(w http.ResponseWriter, r *http.Request) {
	lista, err := api.rep.ListTickets()
	if err != nil {
		http.Error(w, "Erro ao obter a lista de tickets no banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(lista)
	if err != nil {
		http.Error(w, "Erro ao converter a lista para json", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) GetTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	ticket, err := api.rep.GetTicketByID(idInt)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro inexistente", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao obter dados no banco de dados", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("http://users-service:8082/users/%d", ticket.UserID)
	fmt.Printf("INFO: Serviço de tickets fazendo uma requisição interna para: %s\n", url)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Header de autorização ausente", http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		http.Error(w, "Header de autorização mal formatado", http.StatusUnauthorized)
		return
	}

	tokenString := headerParts[1]

	cliente := &http.Client{
		Timeout: 5 * time.Second,
	}

	reqInternal, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar requisição interna: "+err.Error(), http.StatusInternalServerError)
		return
	}

	reqInternal.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	reqInternal.Header.Set("Content-Type", "application/json")

	//Fazendo a requisição
	resposta, err := cliente.Do(reqInternal)
	if err != nil {
		log.Printf("ERRO: Falha ao fazer requisição para o serviço de usuarios: %v", err)
		http.Error(w, "ERRO: Falha ao fazer a requisição para o serviço de usuarios", http.StatusBadRequest)
		return
	}
	defer resposta.Body.Close()

	//Verificando se o status da resposta foi bem sucedido (status code 2xx)
	if resposta.StatusCode != http.StatusOK {
		http.Error(w, "Erro ao consultar serviço de usuarios", http.StatusBadRequest)
		return
	}

	var ticketAuthor model.TicketAuthor
	if err = json.NewDecoder(resposta.Body).Decode(&ticketAuthor); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da resposta do users-service", http.StatusBadRequest)
		return
	}

	ticket.Author = ticketAuthor

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ticket)
	if err != nil {
		http.Error(w, "Erro ao converter ticket para json", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) GetMyTicketsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.UserIDKey).(int64)

	lista, err := api.rep.GetTicketByUser(int(id))
	if err != nil {
		http.Error(w, "Erro ao consultar tickets do usuario", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Erro ao codificar a lista de tarefas", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	idReq := r.Context().Value(middleware.UserIDKey)

	var ticketReq model.UpdateTicketPayload
	err = json.NewDecoder(r.Body).Decode(&ticketReq)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	var ticketOg model.Ticket
	ticketOg, err = api.rep.GetTicketByID(idInt)
	if err != nil {
		http.Error(w, "Erro ao obter o ticket no banco de dados", http.StatusInternalServerError)
		return
	}

	if idReq != ticketOg.UserID {
		http.Error(w, "Permissão não concedida", http.StatusForbidden)
	}

	ticketOg.Titulo = ticketReq.Titulo
	ticketOg.Descricao = ticketReq.Descricao
	ticketOg.Prioridade = ticketReq.Prioridade
	ticketOg.Anexos = ticketReq.Anexos
	ticketOg.CategoriaID = ticketReq.CategoriaID

	err = api.rep.UpdateTicket(idInt, ticketOg)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao modificar registro no banco de dados", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) UpdateTicketStatusHandler(w http.ResponseWriter, r *http.Request) {
	type updateStatusRequest struct {
		ID     int64  `json:"id"`
		Status string `json:"status"`
	}

	var statusReq updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&statusReq); err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}

	ticketOg, err := api.rep.GetTicketByID(int(statusReq.ID))
	if err != nil {
		http.Error(w, "Erro ao consultar o ticket original", http.StatusBadRequest)
		return
	}

	idReq := r.Context().Value(middleware.UserIDKey)

	if idReq != ticketOg.UserID {
		http.Error(w, "Permissão não concedida", http.StatusForbidden)
		return
	}

	ticketOg.Status = statusReq.Status

	if err := api.rep.UpdateTicket(int(ticketOg.ID), ticketOg); err != nil {
		http.Error(w, "Erro ao atualizar informacoes no banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	ticket, err := api.rep.GetTicketByID(idInt)
	if err != nil {
		http.Error(w, "Ticket não encontrado no banco de dados", http.StatusBadRequest)
		return
	}

	idReq := r.Context().Value(middleware.UserIDKey).(int64)

	if ticket.UserID != idReq {
		http.Error(w, "Permissão não concedida", http.StatusForbidden)
		return
	}

	err = api.rep.DeleteTicket(idInt)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao remover registro do banco de dados", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *ApiServer) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do ticket inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	var comentario model.Comentario
	if err = json.NewDecoder(r.Body).Decode(&comentario); err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}

	idUser, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Não foi possivel extrair o ID do usuario do token", http.StatusInternalServerError)
		return
	}

	comentario.UserID = idUser
	comentario.TicketID = int64(idInt)

	idComent, err := api.rep.CreateComment(comentario)
	if err != nil {
		http.Error(w, "Erro ao adicionar o comentario no banco de dados", http.StatusBadRequest)
		return
	}
	comentario.ID = idComent

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comentario)
}

func (api *ApiServer) ListCommentsByTicketHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do ticket inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	lista, err := api.rep.ListCommentsByTicketID(id)
	if err != nil {
		http.Error(w, "Erro ao consultar o BD", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(lista); err != nil {
		// Se a lista foi obtida mas a codificação falha, é um erro do servidor.
		http.Error(w, "Erro ao codificar a resposta em JSON", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) ListCommentsByUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do usuário inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	lista, err := api.rep.ListCommentsByUserID(id)
	if err != nil {
		http.Error(w, "Erro ao consultar o BD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Erro ao codificar a resposta em JSON", http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do comentário inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	idUser, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Erro ao obter o ID do usuario do cabeçalho da requisição", http.StatusUnauthorized)
		return
	}

	var comment model.Comentario
	comment.UserID = idUser

	if err = json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = api.rep.UpdateComment(id, comment); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Comentário não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao atualizar o comentario no banco de dados", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do comentário inválido, deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	err = api.rep.DeleteComment(id)
	if err == pgx.ErrNoRows {
		http.Error(w, "Comentário não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao remover o comentário do banco de dados", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
