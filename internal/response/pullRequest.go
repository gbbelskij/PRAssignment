package response

import (
	"PRAssignment/internal/domain"
	"time"
)

type PullRequestCreateResponse struct {
	PR PullRequestCreate `json:"pr"`
}

type PullRequestCreate struct {
	PullRequestId     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorId          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
}

type PullRequestMergeResponse struct {
	PR PullRequestMerge `json:"pr"`
}

type PullRequestMerge struct {
	PullRequestId     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorId          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
	MergedAt          time.Time                `json:"merged_at"`
}

type PullRequestReassignResponse struct {
	PR         PullRequestReassign `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type PullRequestReassign struct {
	PullRequestId     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorId          string                   `json:"author_id"`
	Status            domain.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
}

func MakePullRequestCreateResponse(
	pullRequestId string,
	pullRequestName string,
	authorId string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
) PullRequestCreateResponse {
	return PullRequestCreateResponse{
		PR: PullRequestCreate{
			PullRequestId:     pullRequestId,
			PullRequestName:   pullRequestName,
			AuthorId:          authorId,
			Status:            status,
			AssignedReviewers: assignedReviewers,
		},
	}
}

func MakePullRequestMergeResponse(
	pullRequestId string,
	pullRequestName string,
	authorId string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
	mergedAt time.Time,
) PullRequestMergeResponse {
	return PullRequestMergeResponse{
		PR: PullRequestMerge{
			PullRequestId:     pullRequestId,
			PullRequestName:   pullRequestName,
			AuthorId:          authorId,
			Status:            status,
			AssignedReviewers: assignedReviewers,
			MergedAt:          mergedAt,
		},
	}
}

func MakePullRequestReassignResponse(
	pullRequestId string,
	pullRequestName string,
	authorId string,
	status domain.PullRequestStatus,
	assignedReviewers []string,
	replacedBy string,
) PullRequestReassignResponse {
	return PullRequestReassignResponse{
		PR: PullRequestReassign{
			PullRequestId:     pullRequestId,
			PullRequestName:   pullRequestName,
			AuthorId:          authorId,
			Status:            status,
			AssignedReviewers: assignedReviewers,
		},
		ReplacedBy: replacedBy,
	}
}
