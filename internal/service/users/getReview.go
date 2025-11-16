package userService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type PullRequestStorage interface {
	GetPullRequestIdsByUserId(ctx context.Context, userID string) ([]string, error)
	GetPullRequestsByIds(ctx context.Context, ids []string) ([]domain.PullRequest, error)
}

type GetReviewService struct {
	storage PullRequestStorage
}

func NewGetReviewService(storage PullRequestStorage) *GetReviewService {
	return &GetReviewService{storage: storage}
}

func (s *GetReviewService) GetUserPullRequests(ctx context.Context, userID string) (*response.UsersGetReviewResponse, error) {
	const op = "service.GetReviewService.GetUserPullRequests"

	pullRequestIds, err := s.storage.GetPullRequestIdsByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(pullRequestIds) == 0 {
		return &response.UsersGetReviewResponse{UserID: userID, PullRequests: []response.PullRequestShort{}}, nil
	}

	pullRequests, err := s.storage.GetPullRequestsByIds(ctx, pullRequestIds)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp := response.MakeUsersGetReviewResponse(userID, pullRequests)

	return &resp, nil
}
