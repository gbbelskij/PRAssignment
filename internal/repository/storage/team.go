package storage

import (
	"PRAssignment/internal/domain"
	"context"
	"fmt"
)

func (s *Storage) SaveTeam(ctx context.Context, team domain.Team) error {
	const op = "repository.storage.team.SaveTeam"

	_, err := s.conn.Exec(
		ctx,
		`INSERT INTO teams(team_name) VALUES($1)`,
		team.TeamName,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
