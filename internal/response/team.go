package response

type TeamGetResponse struct {
	TeamName string       `json:"team_name" env-required:"true"`
	Members  []TeamMember `json:"members" env-required:"true"`
}

type TeamMember struct {
	UserID   string `json:"user_id" env-required:"true"`
	Username string `json:"username" env-required:"true"`
	IsActive bool   `json:"is_active" env-required:"true"`
}
