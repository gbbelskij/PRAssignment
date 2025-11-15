package request

type PullRequestCreateRequest struct {
	PullRequestId   string `json:"pull_request_id" env-required:"true"`
	PullRequestName string `json:"pull_request_name" env-required:"true"`
	AuthorId        string `json:"author_id" env-required:"true"`
}

type PullRequestMergeRequest struct {
	PullRequestId string `json:"pull_request_id" env-required:"true"`
}

type PullRequestReassignRequest struct {
	PullRequestId string `json:"pull_request_id" env-required:"true"`
	OldUserId     string `json:"old_user_id" env-required:"true"`
}
