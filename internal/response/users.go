package response

import "PRAssignment/internal/domain"

type UserSetIsActiveResponse struct {
	User UserSetIsActive `json:"user"`
}

type UserSetIsActive struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UsersGetReviewResponse struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PullRequestShort struct {
	PullRequestId   string                   `json:"pull_request_id"`
	PullRequestName string                   `json:"pull_request_name"`
	AuthorId        string                   `json:"author_id"`
	Status          domain.PullRequestStatus `json:"status"`
}

func MakeUserSetIsActiveResponse(
	userId string,
	username string,
	teamName string,
	isActive bool,
) UserSetIsActiveResponse {
	return UserSetIsActiveResponse{
		User: UserSetIsActive{
			UserId:   userId,
			Username: username,
			TeamName: teamName,
			IsActive: isActive,
		},
	}
}

func MakeUsersGetReviewResponse(userId string, pullRequests []domain.PullRequest) UsersGetReviewResponse {
	prShorts := make([]PullRequestShort, 0, len(pullRequests))
	for _, pr := range pullRequests {
		prShorts = append(prShorts, PullRequestShort{
			PullRequestId:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorID,
			Status:          pr.Status,
		})
	}
	return UsersGetReviewResponse{
		UserId:       userId,
		PullRequests: prShorts,
	}
}
