package response

import "PRAssignment/internal/domain"

type PullRequestCreateResponse struct {
	PR PullRequest `json:"pr"`
}

type PullRequest struct {
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
		PR: PullRequest{
			PullRequestId:     pullRequestId,
			PullRequestName:   pullRequestName,
			AuthorId:          authorId,
			Status:            status,
			AssignedReviewers: assignedReviewers,
		},
	}
}
