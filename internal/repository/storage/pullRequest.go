package storage

import (
	"PRAssignment/internal/domain"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"context"
	"errors"
	"fmt"
	"strings"

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

		reviewers, err = s.FindReviewersTx(txCtx, tx, teamId, pullRequest.AuthorID, pullRequest.PullRequestID)
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

func (s *Storage) FindReviewersTx(ctx context.Context, tx pgx.Tx, teamId string, userId string, pullRequestId string) ([]string, error) {
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
		var reviewerUserId string
		if err := rows.Scan(&reviewerUserId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviewers = append(reviewers, reviewerUserId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, revId := range reviewers {
		_, err := tx.Exec(ctx,
			`INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ($1, $2)`,
			pullRequestId, revId,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return reviewers, nil
}

func (s *Storage) getPullRequestById(ctx context.Context, pullRequestId string) (*domain.PullRequest, error) {
	const op = "repository.storage.getPullRequestById"

	var pullRequest domain.PullRequest
	err := s.conn.QueryRow(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, updated_at, merged_at, created_at
        FROM pull_requests WHERE pull_request_id = $1`,
		pullRequestId,
	).Scan(
		&pullRequest.PullRequestID,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.UpdatedAt,
		&pullRequest.MergedAt,
		&pullRequest.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pullRequest, nil
}

func (s *Storage) getStatusPullRequestById(ctx context.Context, pullRequestId string) (*domain.PullRequestStatus, error) {
	const op = "repository.storage.getPullRequestById"

	var status domain.PullRequestStatus
	err := s.conn.QueryRow(ctx,
		`SELECT status
        FROM pull_requests WHERE pull_request_id = $1`,
		pullRequestId,
	).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &status, nil
}

func (s *Storage) UpdatePullRequest(ctx context.Context, pullRequestId string) (*domain.PullRequest, error) {
	const op = "repository.storage.UpdatePullRequest"

	pr, err := s.getPullRequestById(ctx, pullRequestId)
	if err != nil {
		return nil, err
	}
	if pr.Status == domain.PullRequestStatusMerged {
		return pr, nil
	}

	_, err = s.conn.Exec(ctx,
		`UPDATE pull_requests 
        SET status = 'MERGED', 
            merged_at = NOW(),
            updated_at = NOW()
        WHERE pull_request_id = $1`,
		pullRequestId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.getPullRequestById(ctx, pullRequestId)
}

func (s *Storage) GetPullRequestReviewers(ctx context.Context, pullRequestId string) ([]string, error) {
	const op = "repository.storage.GetPullRequestReviewers"

	var reviewers []string
	rows, err := s.conn.Query(ctx,
		`SELECT user_id FROM pull_request_reviewers WHERE pull_request_id = $1`,
		pullRequestId,
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

func (s *Storage) ReassignReviewer(ctx context.Context, pullRequestId string, oldUserId string) (string, error) {
	const op = "repository.storage.ReassignReviewer"

	var newReviewerId string
	err := s.conn.QueryRow(ctx,
		`SELECT user_id FROM team_members
		WHERE team_id = (SELECT team_id FROM team_members WHERE user_id = $1 LIMIT 1)
		AND user_id <> $1
		AND is_active = true
		LIMIT 1`,
		oldUserId,
	).Scan(&newReviewerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", customErrors.ErrNoCandidate
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx,
		`DELETE FROM pull_request_reviewers WHERE user_id = $1`,
		oldUserId,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx,
		`INSERT INTO pull_request_reviewers(pull_request_id, user_id) VALUES($1, $2)`,
		pullRequestId, newReviewerId,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return newReviewerId, nil
}

func (s *Storage) UpdateReviewer(ctx context.Context, pullRequestId string, oldUserId string) (*domain.PullRequest, string, error) {
	const op = "repository.storage.UpdateReviewer"

	status, err := s.getStatusPullRequestById(ctx, pullRequestId)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if *status == domain.PullRequestStatusMerged {
		return nil, "", customErrors.ErrPrMerged
	}

	reviewers, err := s.GetPullRequestReviewers(ctx, pullRequestId)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	isReviewer := false
	for _, r := range reviewers {
		if r == oldUserId {
			isReviewer = true
			break
		}
	}

	if !isReviewer {
		return nil, "", customErrors.ErrNotAssigned
	}

	newReviewerId, err := s.ReassignReviewer(ctx, pullRequestId, oldUserId)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	pullRequest, err := s.getPullRequestById(ctx, pullRequestId)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return pullRequest, newReviewerId, nil
}

func (s *Storage) GetPullRequests(ctx context.Context, userId string) ([]domain.PullRequest, error) {
	const op = "repository.storage.GetPullRequests"

	rows, err := s.conn.Query(ctx,
		`SELECT pull_request_id FROM pull_request_reviewers WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var pullRequestIds []string
	for rows.Next() {
		var prId string
		if err := rows.Scan(&prId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		pullRequestIds = append(pullRequestIds, prId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(pullRequestIds) == 0 {
		return []domain.PullRequest{}, nil
	}

	query, args, err := buildInQuery(
		`SELECT pull_request_id, pull_request_name, author_id, status, updated_at, merged_at, created_at
		FROM pull_requests
		WHERE pull_request_id IN (`,
		pullRequestIds,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	prRows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer prRows.Close()

	var pullRequests []domain.PullRequest
	for prRows.Next() {
		var pr domain.PullRequest
		if err := prRows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&pr.UpdatedAt,
			&pr.MergedAt,
			&pr.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		pullRequests = append(pullRequests, pr)
	}

	if err := prRows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pullRequests, nil
}

func buildInQuery(baseQuery string, values []string) (string, []interface{}, error) {
	if len(values) == 0 {
		return "", nil, fmt.Errorf("no values to build query")
	}
	placeholders := make([]string, len(values))
	args := make([]interface{}, len(values))
	for i, v := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = v
	}
	query := baseQuery + strings.Join(placeholders, ",") + ")"
	return query, args, nil
}
