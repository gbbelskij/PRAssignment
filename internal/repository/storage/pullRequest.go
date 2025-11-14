package storage

import (
	"PRAssignment/internal/domain"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddPullRequest(ctx context.Context, pullRequest *domain.PullRequest) ([]string, error) {
	const op = "repository.storage.AddPullRequest"

	var reviewers []string
	err := s.txManager.WithTx(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var teamId string
		err := tx.QueryRow(txCtx,
			`SELECT team_id FROM team_members WHERE user_id = $1`,
			pullRequest.AuthorID,
		).Scan(&teamId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return customErrors.ErrNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = tx.Exec(txCtx,
			`INSERT INTO pull_requests(pull_request_id, pull_request_name, author_id, status)
            VALUES($1, $2, $3, $4)`,
			pullRequest.PullRequestID, pullRequest.PullRequestName, pullRequest.AuthorID, pullRequest.Status,
		)
		if err != nil {
			if customErrors.IsUniqueViolation(err) {
				return customErrors.ErrUniqueViolation
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		reviewers, err = s.FindReviewersTx(txCtx, tx, teamId, pullRequest.AuthorID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if len(reviewers) == 0 {
			return customErrors.ErrNoCandidate
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}
	return reviewers, nil
}

func (s *Storage) FindReviewersTx(ctx context.Context, tx pgx.Tx, teamId string, userId string) ([]string, error) {
	const op = "repository.storage.FindReviewersTx"

	var reviewers []string
	rows, err := tx.Query(ctx,
		`SELECT user_id FROM team_members WHERE team_id = $1 AND user_id <> $2 AND is_active = true
        LIMIT 2`,
		teamId, userId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviewers = append(reviewers, userId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviewers, nil
}
