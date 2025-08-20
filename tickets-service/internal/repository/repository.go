package repository

import (
	"context"
	"errors"
	"helpdesk/tickets-service/internal/model"
	"log"

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
	if err := s.db.QueryRow(context.Background(), "INSERT INTO tickets (titulo, descricao, status, diagnostico, solucao, prioridade, data_abertura, data_fechamento, data_atualizacao, anexos, tags, categoria_id, responsavel_id, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) returning id", ticket.Titulo, ticket.Descricao, ticket.Status, ticket.Diagnostico, ticket.Solucao, ticket.Prioridade, ticket.DataAbertura, ticket.DataFechamento, ticket.DataAtualizacao, ticket.Anexos, ticket.Tags, ticket.CategoriaID, ticket.ResponsavelID, ticket.UserID).Scan(&ticket.ID); err != nil {
		go func() {
			log.Printf("Erro ao adicionar ticket no banco de dados: %v", err)
		}()
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
		if err := rows.Scan(&ticket.ID, &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID); err != nil {
			return nil, err
		}
		lista = append(lista, ticket)
	}

	return lista, nil
}

func (s *Repository) GetTicketByID(id int) (model.Ticket, error) {
	var ticket model.Ticket
	if err := s.db.QueryRow(context.Background(), "SELECT * FROM tickets WHERE id=$1", id).Scan(&ticket.ID, &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID); err != nil {
		return ticket, err
	}

	return ticket, nil
}

func (s *Repository) GetTicketByUser(id int) ([]model.Ticket, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM tickets WHERE user_id=$1", id)
	if err != nil {
		return []model.Ticket{}, err
	}

	var ticket model.Ticket
	var lista []model.Ticket

	for rows.Next() {
		if err = rows.Scan(&ticket.ID, &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID); err != nil {
			return []model.Ticket{}, err
		}
		lista = append(lista, ticket)
	}

	return lista, nil
}

func (s *Repository) UpdateTicket(id int, ticket model.Ticket) error {
	_, err := s.db.Exec(context.Background(), "UPDATE tickets SET titulo=$1, descricao=$2, status=$3, diagnostico=$4, solucao=$5, prioridade=$6, data_abertura=$7, data_fechamento=$8, data_atualizacao=$9, anexos=$10, tags=$11, categoria_id=$12, responsavel_id=$13, user_id=$14 WHERE id=$15", &ticket.Titulo, &ticket.Descricao, &ticket.Status, &ticket.Diagnostico, &ticket.Solucao, &ticket.Prioridade, &ticket.DataAbertura, &ticket.DataFechamento, &ticket.DataAtualizacao, &ticket.Anexos, &ticket.Tags, &ticket.CategoriaID, &ticket.ResponsavelID, &ticket.UserID, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return pgx.ErrNoRows
	} else if err != nil {
		return err
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

func (s *Repository) CreateComment(comment model.Comentario) (int64, error) {
	if err := s.db.QueryRow(context.Background(), "INSERT INTO comentarios (descricao, data, user_id, ticket_id) VALUES ($1, $2, $3, $4) returning id", comment.Descricao, comment.Data, comment.UserID, comment.TicketID).Scan(&comment.ID); err != nil {
		go func() {
			log.Printf("Erro ao adicionar comentario no banco de dados: %v", err)
		}()
		return 0, err
	}

	return comment.ID, nil
}

func (s *Repository) GetCommentByID(id int) (model.Comentario, error) {
	var comentario model.Comentario

	if err := s.db.QueryRow(context.Background(), "SELECT * FROM comentarios WHERE id=$1", id).Scan(&comentario.ID, &comentario.Descricao, &comentario.Data, &comentario.UserID, &comentario.TicketID); err != nil {
		return model.Comentario{}, err
	}

	return comentario, nil
}

func (s *Repository) ListCommentsByTicketID(id int) ([]model.Comentario, error) {
	var lista []model.Comentario
	var comentario model.Comentario

	rows, err := s.db.Query(context.Background(), "SELECT * FROM comentarios WHERE ticket_id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&comentario.ID, &comentario.Descricao, &comentario.Data, &comentario.UserID, &comentario.TicketID); err != nil {
			return nil, err
		}
		lista = append(lista, comentario)
	}

	return lista, nil
}

func (s *Repository) ListCommentsByUserID(id int) ([]model.Comentario, error) {
	var lista []model.Comentario
	var comentario model.Comentario

	rows, err := s.db.Query(context.Background(), "SELECT * FROM comentarios WHERE user_id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&comentario.ID, &comentario.Descricao, &comentario.Data, &comentario.UserID, &comentario.TicketID); err != nil {
			return nil, err
		}
		lista = append(lista, comentario)
	}

	return lista, nil
}

func (s *Repository) UpdateComment(id int, comment model.Comentario) error {
	_, err := s.db.Exec(context.Background(), "UPDATE comentarios SET descricao=$1, data=$2, user_id=$3 WHERE id=$4", &comment.Descricao, &comment.Data, &comment.UserID, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func (s *Repository) DeleteComment(id int) error {
	_, err := s.db.Exec(context.Background(), "DELETE FROM comentarios WHERE id=$1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return pgx.ErrNoRows
	} else if err != nil {
		return err
	}
	return nil
}
