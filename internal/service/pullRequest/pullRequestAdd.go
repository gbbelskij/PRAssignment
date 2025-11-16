package pullRequestService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type PullRequestStorage interface {
	AddPullRequest(ctx context.Context, pullRequest *domain.PullRequest) ([]string, error)
}

type PullRequestCreateService struct {
	storage PullRequestStorage
}

func NewPullRequestCreateService(storage PullRequestStorage) *PullRequestCreateService {
	return &PullRequestCreateService{storage: storage}
}

func (s *PullRequestCreateService) CreatePullRequest(ctx context.Context, req *request.PullRequestCreateRequest) (*response.PullRequestCreateResponse, error) {
	pullRequest := PullRequestFromRequest(*req)

	reviewers, err := s.storage.AddPullRequest(ctx, pullRequest)
	if err != nil {
		return nil, fmt.Errorf("add pull request: %w", err)
	}

	resp := response.MakePullRequestCreateResponse(
		req.PullRequestID,
		req.PullRequestName,
		req.AuthorID,
		domain.PullRequestStatusOpen,
		reviewers,
	)

	return &resp, nil
}

func PullRequestFromRequest(prRequest request.PullRequestCreateRequest) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   prRequest.PullRequestID,
		PullRequestName: prRequest.PullRequestName,
		AuthorID:        prRequest.AuthorID,
		Status:          domain.PullRequestStatusOpen,
	}
}
