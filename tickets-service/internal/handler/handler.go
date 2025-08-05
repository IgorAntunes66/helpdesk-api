package handler

import (
	"encoding/json"
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
		http.Error(w, "Erro ao decodigicar o corpo da requisição", http.StatusInternalServerError)
		return
	}

	id, err := api.rep.CreateTicket(ticket)
	if err != nil {
		http.Error(w, "Erro ao adicionar o ticket no banco de dados", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(lista)
	if err != nil {
		http.Error(w, "Erro ao converter a lista para json", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) GetTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter id para inteiro", http.StatusInternalServerError)
		return
	}

	ticket, err := api.rep.GetTicketByID(idInt)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro inexistente", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao obter dados no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ticket)
	if err != nil {
		http.Error(w, "Erro ao converter ticket para json", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao converter id para inteiro", http.StatusInternalServerError)
		return
	}

	var ticket model.Ticket
	err = json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusInternalServerError)
		return
	}

	err = api.rep.UpdateTicket(idInt, ticket)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao modificar registro no banco de dados", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *ApiServer) DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao obter ID da requisição", http.StatusInternalServerError)
		return
	}

	err = api.rep.DeleteTicket(idInt)
	if err == pgx.ErrNoRows {
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro ao remover registro do banco de dados", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
