package domain

import "time"

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "open"
	PullRequestStatusClosed PullRequestStatus = "closed"
)

type PullRequest struct {
	PullRequestID   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	AuthorID        string            `json:"author_id"`
	Status          PullRequestStatus `json:"status"`
	UpdatedAt       time.Time         `json:"updated_at"`
	MergedAt        time.Time         `json:"merged_at"`
	CreatedAt       time.Time         `json:"created_at"`
}

type PullRequestReviewer struct {
	PullRequestID string `json:"pull_request_id"`
	UserID        string `json:"user_id"`
}
