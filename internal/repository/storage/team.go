package storage

import (
	"PRAssignment/internal/domain"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"context"
	"fmt"
)

func (s *Storage) SaveTeam(ctx context.Context, team *domain.Team) (string, error) {
	const op = "repository.storage.team.SaveTeam"

	var teamID string
	err := s.conn.QueryRow(
		ctx,
		`INSERT INTO teams(team_name) VALUES($1)
        RETURNING team_id`,
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

func (s *Storage) SaveTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	const op = "repository.storage.team.SaveTeamMember"

	_, err := s.conn.Exec(
		ctx,
		`INSERT INTO team_members(team_id, user_id, username, is_active) VALUES($1, $2, $3, $4)`,
		teamMember.TeamID, teamMember.UserID, teamMember.Username, teamMember.IsActive,
	)
	if err != nil {
		if customErrors.IsUniqueViolation(err) {
			return customErrors.ErrUniqueViolation
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveTeamMembersBatch(ctx context.Context, teamMembers []domain.TeamMember) error {
	const op = "repository.storage.team.SaveTeamMembersBatch"

	for _, teamMember := range teamMembers {
		err := s.SaveTeamMember(ctx, &teamMember)
		if err != nil {
			if customErrors.IsUniqueViolation(err) {
				return customErrors.ErrUniqueViolation
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
