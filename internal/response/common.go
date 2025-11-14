package response

type ErrorCode string

const (
	ErrCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists    ErrorCode = "PR_EXISTS"
	ErrCodePRMerged    ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest  ErrorCode = "BAD_REQUEST"
)

type ErrorResponse struct {
	Error Error `json:"error" env-required:"true"`
}

type Error struct {
	Code    ErrorCode `json:"code" env-required:"true"`
	Message string    `json:"message" env-required:"true"`
}

func MakeError(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}
