package domain

import "time"

type Team struct {
	TeamID    string    `json:"team_id"`
	TeamName  string    `json:"team_name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	TeamID   string `json:"team_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
