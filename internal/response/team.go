package response

type ErrorCode string

const (
	ErrCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists    ErrorCode = "PR_EXISTS"
	ErrCodePRMerged    ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrBadRequest      ErrorCode = "BAD_REQUEST"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type TeamGetResponse struct {
	TeamName string       `json:"team_name" env-required:"true"`
	Members  []TeamMember `json:"members" env-required:"true"`
}

type TeamMember struct {
	UserId   string `json:"user_id" env-required:"true"`
	Username string `json:"username" env-required:"true"`
	IsActive bool   `json:"is_active" env-required:"true"`
}
