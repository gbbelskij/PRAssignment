package service

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/request"
)

func PullRequestFromRequest(prRequest request.PullRequestCreateRequest) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   prRequest.PullRequestId,
		PullRequestName: prRequest.PullRequestName,
		AuthorID:        prRequest.AuthorId,
		Status:          domain.PullRequestStatusOpen,
	}
}
