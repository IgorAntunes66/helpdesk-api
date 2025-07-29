package repository

import (
	"context"
	"helpdesk/users-service/internal/model"

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

func (s *Repository) CreateUser(user model.User) (int, error) {
	err := s.db.QueryRow(context.Background(), "INSERT INTO users (nome, senha, funcao, telefone, cpfCnpj) VALUES ($1, $2, $3, $4, $5) returning id", user.Nome, user.Senha, user.Funcao, user.Telefone, user.CpfCnpj).Scan(&user.ID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *Repository) FindAllUsers() ([]model.User, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.User
	var u model.User

	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Nome, &u.Senha, &u.Funcao, &u.Telefone, &u.CpfCnpj); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

func (s *Repository) FindUserByID(id int) (model.User, error) {
	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id=$1", id)

	var u model.User

	if err := row.Scan(&u.ID, &u.Nome, &u.Senha, &u.Funcao, &u.Telefone, &u.CpfCnpj); err != nil {
		return u, err
	}

	return u, nil
}

func (s *Repository) UpdateUser(id int, user model.User) error {
	row, err := s.db.Exec(context.Background(), "UPDATE users SET nome=$1, senha=$2, funcao=$3, telefone=$4, cpfCnpj=$5 WHERE id=$6", user.Nome, user.Senha, user.Funcao, user.Telefone, user.CpfCnpj, id)
	if err != nil {
		return err
	}

	if row.RowsAffected() != 1 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Repository) DeleteUser(id int) error {
	row, err := s.db.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	if row.RowsAffected() != 1 {
		return pgx.ErrNoRows
	}

	return nil
}
