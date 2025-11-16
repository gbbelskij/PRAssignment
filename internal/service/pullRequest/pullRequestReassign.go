package pullRequestService

import (
	"PRAssignment/internal/domain"
	customerrors "PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type PullRequestReassignStorage interface {
	GetPullRequestById(ctx context.Context, pullRequestID string) (*domain.PullRequest, error)
	GetPullRequestStatus(ctx context.Context, pullRequestID string) (*domain.PullRequestStatus, error)
	GetPullRequestReviewers(ctx context.Context, pullRequestID string) ([]string, error)
	ReassignReviewerInDb(ctx context.Context, pullRequestID string, oldUserID string) (string, error)
}

type PullRequestReassignService struct {
	storage PullRequestReassignStorage
}

func NewPullRequestReassignService(storage PullRequestReassignStorage) *PullRequestReassignService {
	return &PullRequestReassignService{storage: storage}
}

func (s *PullRequestReassignService) ReassignReviewer(ctx context.Context, pullRequestID string, oldUserID string) (*response.PullRequestReassignResponse, error) {
	status, err := s.storage.GetPullRequestStatus(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get pull request status: %w", err)
	}

	if *status == domain.PullRequestStatusMerged {
		return nil, customerrors.ErrPrMerged
	}

	reviewers, err := s.storage.GetPullRequestReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get reviewers: %w", err)
	}

	isReviewer := false
	for _, r := range reviewers {
		if r == oldUserID {
			isReviewer = true
			break
		}
	}

	if !isReviewer {
		return nil, customerrors.ErrNotAssigned
	}

	newReviewerID, err := s.storage.ReassignReviewerInDb(ctx, pullRequestID, oldUserID)
	if err != nil {
		return nil, fmt.Errorf("reassign reviewer: %w", err)
	}

	pr, err := s.storage.GetPullRequestById(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get updated pull request: %w", err)
	}

	updatedReviewers, err := s.storage.GetPullRequestReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get updated reviewers: %w", err)
	}

	resp := response.MakePullRequestReassignResponse(
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		updatedReviewers,
		newReviewerID,
	)

	return &resp, nil
}
