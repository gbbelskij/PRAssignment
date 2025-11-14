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

func (s *Storage) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
	}()

	err = fn(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) Close() {
	s.conn.Close()
}
