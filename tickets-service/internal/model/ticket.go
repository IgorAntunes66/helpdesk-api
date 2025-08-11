package model

import "time"

type Ticket struct {
	ID              int64        `json:"id"`
	Titulo          string       `json:"titulo"`
	Descricao       string       `json:"descricao"`
	Status          string       `json:"status"`
	Diagnostico     string       `json:"diagnostico"`
	Solucao         string       `json:"solucao"`
	Prioridade      string       `json:"prioridade"`
	DataAbertura    time.Time    `json:"data_abertura"`
	DataFechamento  time.Time    `json:"data_fechamento"`
	DataAtualizacao time.Time    `json:"data_atualizacao"`
	Anexos          []string     `json:"anexos"`
	Tags            []string     `json:"tags"`
	Historico       []Comentario `json:"historico"`
	CategoriaID     int64        `json:"categoria_id"`
	ResponsavelID   int64        `json:"responsavel_id"`
	UserID          int64        `json:"user_id"`
}

type Comentario struct {
	ID        int64     `json:"id"`
	Descricao string    `json:"descricao"`
	Data      time.Time `json:"data"`
	UserID    int64     `json:"user_id"`
	TicketID  int64     `json:"ticket_id"`
}

type UpdateTicketPayload struct {
	Titulo      string   `json:"titulo"`
	Descricao   string   `json:"descricao"`
	Prioridade  string   `json:"prioridade"`
	Anexos      []string `json:"anexos"`
	CategoriaID int64    `json:"categoria_id"`
}
