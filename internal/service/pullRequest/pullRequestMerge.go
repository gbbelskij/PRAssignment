package pullRequestService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type PullRequestMergeStorage interface {
	GetPullRequestById(ctx context.Context, pullRequestID string) (*domain.PullRequest, error)
	UpdatePullRequestStatus(ctx context.Context, pullRequestID string) error
	GetPullRequestReviewers(ctx context.Context, pullRequestID string) ([]string, error)
}

type PullRequestMergeService struct {
	storage PullRequestMergeStorage
}

func NewPullRequestMergeService(storage PullRequestMergeStorage) *PullRequestMergeService {
	return &PullRequestMergeService{storage: storage}
}

func (s *PullRequestMergeService) MergePullRequest(ctx context.Context, pullRequestID string) (*response.PullRequestMergeResponse, error) {
	pr, err := s.storage.GetPullRequestById(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get pull request: %w", err)
	}

	if pr.Status != domain.PullRequestStatusMerged {
		err = s.storage.UpdatePullRequestStatus(ctx, pullRequestID)
		if err != nil {
			return nil, fmt.Errorf("update pull request status: %w", err)
		}
		pr, err = s.storage.GetPullRequestById(ctx, pullRequestID)
		if err != nil {
			return nil, fmt.Errorf("get updated pull request: %w", err)
		}
	}

	reviewers, err := s.storage.GetPullRequestReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get reviewers: %w", err)
	}

	resp := response.MakePullRequestMergeResponse(
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		reviewers,
		pr.MergedAt,
	)

	return &resp, nil
}
