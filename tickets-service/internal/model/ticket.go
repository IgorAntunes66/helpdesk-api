package model

type Ticket struct {
	ID          int    `json:"id"`
	Titulo      string `json:"titulo"`
	Descricao   string `json:"descricao"`
	Status      string `json:"status"`
	Diagnostico string `json:"diagnostico"`
	Solucao     string `json:"solucao"`
	Prioridade  string `json:"prioridade"`
	UserID      int    `json:"user_id"`
}