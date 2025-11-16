package storage

import (
	"PRAssignment/internal/domain"
	customerrors "PRAssignment/internal/repository/customErrors"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) UserExists(ctx context.Context, tx pgx.Tx, userID string) (bool, error) {
	const op = "repository.storage.UserExists"

	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM team_members WHERE user_id = $1)`,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *Storage) UpdateUser(ctx context.Context, userID string, isActive bool) (*domain.TeamMember, error) {
	const op = "repository.storage.UpdateUser"

	cmdTag, err := s.conn.Exec(ctx,
		`UPDATE team_members SET is_active = $1 WHERE user_id = $2`,
		isActive, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, customerrors.ErrNotFound
	}

	var teamMember domain.TeamMember
	err = s.conn.QueryRow(ctx,
		`SELECT user_id, team_id, username, is_active FROM team_members WHERE user_id = $1`,
		userID,
	).Scan(&teamMember.UserID, &teamMember.TeamID, &teamMember.Username, &teamMember.IsActive)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &teamMember, nil
}
