package response

import (
	"PRAssignment/internal/domain"
	"time"
)

type PullRequestCreateResponse struct {
	PR PullRequestCreate `json:"pr"`
}

type PullRequestCreate struct {
	PullRequestID     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorID          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
}

type PullRequestMergeResponse struct {
	PR PullRequestMerge `json:"pr"`
}

type PullRequestMerge struct {
	PullRequestID     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorID          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
	MergedAt          time.Time                `json:"merged_at"`
}

type PullRequestReassignResponse struct {
	PR         PullRequestReassign `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type PullRequestReassign struct {
	PullRequestID     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorID          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
}

func MakePullRequestCreateResponse(
	pullRequestID string,
	pullRequestName string,
	authorID string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
) PullRequestCreateResponse {
	return PullRequestCreateResponse{
		PR: PullRequestCreate{
			PullRequestID:     pullRequestID,
			PullRequestName:   pullRequestName,
			AuthorID:          authorID,
			Status:            status,
			AssignedReviewers: assignedReviewers,
		},
	}
}

func MakePullRequestMergeResponse(
	pullRequestID string,
	pullRequestName string,
	authorID string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
	mergedAt time.Time,
) PullRequestMergeResponse {
	return PullRequestMergeResponse{
		PR: PullRequestMerge{
			PullRequestID:     pullRequestID,
			PullRequestName:   pullRequestName,
			AuthorID:          authorID,
			Status:            status,
			AssignedReviewers: assignedReviewers,
			MergedAt:          mergedAt,
		},
	}
}

func MakePullRequestReassignResponse(
	pullRequestID string,
	pullRequestName string,
	authorID string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
	replacedBy string,
) PullRequestReassignResponse {
	return PullRequestReassignResponse{
		PR: PullRequestReassign{
			PullRequestID:     pullRequestID,
			PullRequestName:   pullRequestName,
			AuthorID:          authorID,
			Status:            status,
			AssignedReviewers: assignedReviewers,
		},
		ReplacedBy: replacedBy,
	}
}
