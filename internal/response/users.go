package response

type UserSetIsActiveResponse struct {
	User User `json:"user"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func MakeUserSetIsActiveResponse(
	userId string,
	username string,
	teamName string,
	isActive bool,
) UserSetIsActiveResponse {
	return UserSetIsActiveResponse{
		User: User{
			UserId:   userId,
			Username: username,
			TeamName: teamName,
			IsActive: isActive,
		},
	}
}
