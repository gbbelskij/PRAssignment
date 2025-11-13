package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(ctx context.Context) (*Storage, error) {
	const op = "repository.New"

	connectStr := os.Getenv("POSTGRES_CONN_STRING")
	conn, err := pgxpool.New(ctx, connectStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}
