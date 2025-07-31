package repository

import (
	"context"
	"helpdesk/tickets-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) CreateTicket(ticket model.Ticket) (int64, error) {
	err := s.db.QueryRow(context.Background(), "INSERT INTO tickets (titulo, descricao, status, diagnostico, solucao, prioridade, data_abertura, data_fechamento, data_atualizacao, anexos, tags, historico, categoria_id, responsavel_id, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning id", ticket.Titulo, ticket.Descricao, ticket.Status, ticket.Diagnostico, ticket.Solucao, ticket.Prioridade, ticket.DataAbertura, ticket.DataFechamento, ticket.DataAtualizacao, ticket.Anexos, ticket.Tags, ticket.Historico, ticket.CategoriaID, ticket.ResponsavelID, ticket.UserID).Scan(&ticket.ID)
	if err != nil {
		return 0, err
	}
	return ticket.ID, nil
}

func (s *Repository) ListTickets() ([]model.Ticket, error) {
	var lista []model.Ticket
	var ticket model.Ticket
	rows, err := s.db.Query(context.Background(), "SELECT * FROM tickets")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&ticket.ID, &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.Historico, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID)
		if err != nil {
			return nil, err
		}
		lista = append(lista, ticket)
	}

	return lista, nil
}

func (s *Repository) GetTicketByID(id int) (model.Ticket, error) {
	var ticket model.Ticket
	err := s.db.QueryRow(context.Background(), "SELECT * FROM tickets WHERE id=$1", id).Scan(&ticket.ID, &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.Historico, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID)
	if err != nil {
		return ticket, err
	}

	return ticket, nil
}

func (s *Repository) UpdateTicket(id int, ticket model.Ticket) error {
	row, err := s.db.Exec(context.Background(), "UPDATE tickets SET titulo=$1, descricao=$2, status=$3, diagnostico=$4, solucao=$5, prioridade=$6, data_abertura=$7, data_fechamento=$8, data_atualizacao=$9, anexos=$10, tags=$11, historico=$12, categoria_id=$13, responsavel_id=$14, user_id=$15 WHERE id=$16", &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.Historico, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID, id)
	if err != nil {
		return err
	}

	if row.RowsAffected() != 1 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Repository) DeleteTicket(id int) error {
	row, err := s.db.Exec(context.Background(), "DELETE FROM tickets WHERE id=$1", id)
	if err != nil {
		return err
	}

	if row.RowsAffected() != 1 {
		return pgx.ErrNoRows
	}

	return nil
}
