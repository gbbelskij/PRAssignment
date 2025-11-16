package response

import "PRAssignment/internal/domain"

type UserSetIsActiveResponse struct {
	User UserSetIsActive `json:"user"`
}

type UserSetIsActive struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UsersGetReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PullRequestShort struct {
	PullRequestID   string                   `json:"pull_request_id"`
	PullRequestName string                   `json:"pull_request_name"`
	AuthorID        string                   `json:"author_id"`
	Status          domain.PullRequestStatus `json:"status"`
}

func MakeUserSetIsActiveResponse(
	userID string,
	username string,
	teamName string,
	isActive bool,
) UserSetIsActiveResponse {
	return UserSetIsActiveResponse{
		User: UserSetIsActive{
			UserID:   userID,
			Username: username,
			TeamName: teamName,
			IsActive: isActive,
		},
	}
}

func MakeUsersGetReviewResponse(userID string, pullRequests []domain.PullRequest) UsersGetReviewResponse {
	prShorts := make([]PullRequestShort, 0, len(pullRequests))
	for _, pr := range pullRequests {
		prShorts = append(prShorts, PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}
	return UsersGetReviewResponse{
		UserID:       userID,
		PullRequests: prShorts,
	}
}
