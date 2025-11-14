package storage

import (
	"PRAssignment/internal/domain"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) SaveTeam(ctx context.Context, tx pgx.Tx, team *domain.Team) (string, error) {
	const op = "repository.storage.SaveTeam"
	var teamID string

	err := tx.QueryRow(ctx,
		`INSERT INTO teams(team_name) VALUES($1) RETURNING team_id`,
		team.TeamName,
	).Scan(&teamID)

	if err != nil {
		if customErrors.IsUniqueViolation(err) {
			return "", customErrors.ErrUniqueViolation
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return teamID, nil
}

func (s *Storage) SaveTeamMember(ctx context.Context, tx pgx.Tx, member *domain.TeamMember) error {
	const op = "repository.storage.SaveTeamMember"

	_, err := tx.Exec(ctx,
		`INSERT INTO team_members(user_id, team_id, username, is_active) VALUES($1, $2, $3, $4)`,
		member.UserID, member.TeamID, member.Username, member.IsActive,
	)

	if err != nil {
		if customErrors.IsUniqueViolation(err) {
			return customErrors.ErrUniqueViolation
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateTeamMemberTeamId(ctx context.Context, tx pgx.Tx, userId string, teamId string) error {
	const op = "repository.storage.UpdateTeamMemberTeamId"

	_, err := tx.Exec(ctx,
		`UPDATE team_members SET team_id = $1 WHERE user_id = $2`,
		teamId, userId,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UserExists(ctx context.Context, tx pgx.Tx, userId string) (bool, error) {
	const op = "repository.storage.UserExists"

	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM team_members WHERE user_id = $1)`,
		userId,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *Storage) SaveTeamWithMembers(
	ctx context.Context,
	team *domain.Team,
	members []domain.TeamMember,
) (string, error) {
	const op = "repository.storage.SaveTeamWithMembers"
	var teamID string

	err := s.txManager.WithTx(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var err error

		teamID, err = s.SaveTeam(txCtx, tx, team)
		if err != nil {
			return fmt.Errorf("%s: failed to save team: %w", op, err)
		}

		for _, member := range members {
			member.TeamID = teamID

			exists, err := s.UserExists(txCtx, tx, member.UserID)
			if err != nil {
				return fmt.Errorf("%s: failed to save team member: %w", op, err)
			}

			if exists {
				if err := s.UpdateTeamMemberTeamId(txCtx, tx, member.UserID, teamID); err != nil {
					return fmt.Errorf("%s: failed to update team member: %w", op, err)
				}
			} else {
				if err := s.SaveTeamMember(txCtx, tx, &member); err != nil {
					return fmt.Errorf("%s: failed to save team member: %w", op, err)
				}
			}
		}

		return nil
	})

	return teamID, err
}

func (s *Storage) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	const op = "repository.storage.team.GetTeam"

	var team domain.Team
	err := s.conn.QueryRow(ctx,
		`SELECT team_id, team_name FROM teams WHERE team_name = $1`,
		teamName).Scan(&team.TeamID, &team.TeamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &team, nil
}

func (s *Storage) GetMembers(ctx context.Context, teamId string) ([]domain.TeamMember, error) {
	const op = "repository.storage.team.GetMembers"

	rows, err := s.conn.Query(
		ctx,
		`SELECT user_id, team_id, username, is_active
		FROM team_members
		WHERE team_id = $1`,
		teamId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var member domain.TeamMember
		if err := rows.Scan(&member.UserID, &member.TeamID, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}
