package storage

import (
	"PRAssignment/internal/domain"
	"context"
	"fmt"
)

func (s *Storage) GetStats(ctx context.Context) (*domain.Stats, error) {
	const op = "repository.storage.GetStats"

	var stats domain.Stats

	err := s.conn.QueryRow(ctx,
		`SELECT 
            (SELECT COUNT(DISTINCT user_id) FROM team_members) as total_users,
            (SELECT COUNT(DISTINCT team_id) FROM team_members) as total_teams,
            COUNT(*) as total_prs,
            COUNT(CASE WHEN status = 'MERGED' THEN 1 END) as merged_prs,
            COUNT(CASE WHEN status = 'OPEN' THEN 1 END) as open_prs
        FROM pull_requests`,
	).Scan(&stats.TotalUsers, &stats.TotalTeams, &stats.TotalPRs, &stats.MergedPRs, &stats.OpenPRs)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	reviewerStats, err := s.conn.Query(ctx,
		`SELECT user_id, COUNT(*) as count
        FROM pull_request_reviewers
        GROUP BY user_id
        ORDER BY count DESC
        LIMIT 10`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer reviewerStats.Close()

	stats.TopReviewers = make([]domain.ReviewerStat, 0)
	for reviewerStats.Next() {
		var rev domain.ReviewerStat
		if err := reviewerStats.Scan(&rev.UserID, &rev.ReviewCount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stats.TopReviewers = append(stats.TopReviewers, rev)
	}

	return &stats, nil
}
