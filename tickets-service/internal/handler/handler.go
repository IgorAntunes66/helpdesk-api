package handler

import (
	"encoding/json"
	"helpdesk/pkg/middleware"
	"helpdesk/tickets-service/internal/model"
	"helpdesk/tickets-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ticket)
	if err != nil {
		http.Error(w, "Erro ao codificar o ticket em json", http.StatusInternalServerError)
		return
	}
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

	var ticket model.Ticket
	err = json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	err = api.rep.UpdateTicket(idInt, ticket)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao modificar registro no banco de dados", http.StatusBadRequest)
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

	var comment model.Comentario
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
