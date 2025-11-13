package models

type TeamMember struct {
	UserId   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
