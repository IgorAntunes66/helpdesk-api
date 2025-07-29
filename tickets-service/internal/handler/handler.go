package handler

import (
	"encoding/json"
	"helpdesk/tickets-service/internal/model"
	"helpdesk/tickets-service/internal/repository"
	"net/http"
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
	}

	id, err := api.rep.CreateTicket(ticket)
	if err != nil {
		http.Error(w, "Erro ao adicionar o ticket no banco de dados", http.StatusInternalServerError)
	}
	ticket.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ticket)
	if err != nil {
		http.Error(w, "Erro ao codificar o ticket em json", http.StatusInternalServerError)
	}
}
