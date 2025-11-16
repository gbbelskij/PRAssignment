package domain

type Stats struct {
	TotalUsers   int            `json:"total_users"`
	TotalTeams   int            `json:"total_teams"`
	TotalPRs     int            `json:"total_prs"`
	MergedPRs    int            `json:"merged_prs"`
	OpenPRs      int            `json:"open_prs"`
	TopReviewers []ReviewerStat `json:"top_reviewers"`
}

type ReviewerStat struct {
	UserID      string `json:"user_id"`
	ReviewCount int    `json:"review_count"`
}
