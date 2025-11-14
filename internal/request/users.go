package request

type UserSetIsActiveRequest struct {
	UserId   string `json:"user_id" env-required:"true"`
	IsActive bool   `json:"is_active" env-required:"true"`
}
