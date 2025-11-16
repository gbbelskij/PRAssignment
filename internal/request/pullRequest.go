package request

type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" env-required:"true"`
	PullRequestName string `json:"pull_request_name" env-required:"true"`
	AuthorID        string `json:"author_id" env-required:"true"`
}

type PullRequestMergeRequest struct {
	PullRequestID string `json:"pull_request_id" env-required:"true"`
}

type PullRequestReassignRequest struct {
	PullRequestID string `json:"pull_request_id" env-required:"true"`
	OldUserID     string `json:"old_user_id" env-required:"true"`
}
