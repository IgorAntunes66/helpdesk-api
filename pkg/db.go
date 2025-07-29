package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var chaveDB string = os.Getenv("CHAVEDB")

func ConectaDB() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), chaveDB)
	if err != nil {
		return nil, fmt.Errorf("não foi possivel conectar ao banco de dados: %w", err)
	}
	return pool, nil
}
