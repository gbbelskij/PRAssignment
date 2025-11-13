package domain

import "time"

type User struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
