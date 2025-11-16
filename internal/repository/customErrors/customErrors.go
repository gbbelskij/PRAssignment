package customErrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}

var (
	ErrUniqueViolation = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrNoCandidate     = errors.New("no candidate")
	ErrPrMerged        = errors.New("pr merged")
	ErrNotAssigned     = errors.New("not assigned")
)
