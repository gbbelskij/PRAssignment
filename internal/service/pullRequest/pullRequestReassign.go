package pullRequestService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type PullRequestReassignStorage interface {
	GetPullRequestById(ctx context.Context, pullRequestId string) (*domain.PullRequest, error)
	GetPullRequestStatus(ctx context.Context, pullRequestId string) (*domain.PullRequestStatus, error)
	GetPullRequestReviewers(ctx context.Context, pullRequestId string) ([]string, error)
	ReassignReviewerInDb(ctx context.Context, pullRequestId string, oldUserId string) (string, error)
}

type PullRequestReassignService struct {
	storage PullRequestReassignStorage
}

func NewPullRequestReassignService(storage PullRequestReassignStorage) *PullRequestReassignService {
	return &PullRequestReassignService{storage: storage}
}

func (s *PullRequestReassignService) ReassignReviewer(ctx context.Context, pullRequestId string, oldUserId string) (*response.PullRequestReassignResponse, error) {
	status, err := s.storage.GetPullRequestStatus(ctx, pullRequestId)
	if err != nil {
		return nil, fmt.Errorf("get pull request status: %w", err)
	}

	if *status == domain.PullRequestStatusMerged {
		return nil, customErrors.ErrPrMerged
	}

	reviewers, err := s.storage.GetPullRequestReviewers(ctx, pullRequestId)
	if err != nil {
		return nil, fmt.Errorf("get reviewers: %w", err)
	}

	isReviewer := false
	for _, r := range reviewers {
		if r == oldUserId {
			isReviewer = true
			break
		}
	}

	if !isReviewer {
		return nil, customErrors.ErrNotAssigned
	}

	newReviewerId, err := s.storage.ReassignReviewerInDb(ctx, pullRequestId, oldUserId)
	if err != nil {
		return nil, fmt.Errorf("reassign reviewer: %w", err)
	}

	pr, err := s.storage.GetPullRequestById(ctx, pullRequestId)
	if err != nil {
		return nil, fmt.Errorf("get updated pull request: %w", err)
	}

	updatedReviewers, err := s.storage.GetPullRequestReviewers(ctx, pullRequestId)
	if err != nil {
		return nil, fmt.Errorf("get updated reviewers: %w", err)
	}

	resp := response.MakePullRequestReassignResponse(
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		updatedReviewers,
		newReviewerId,
	)

	return &resp, nil
}
