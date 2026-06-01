package domain

import "time"

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	Members   []string  `json:"members"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TeamMember struct {
	TeamID   string
	UserID   string
	JoinedAt time.Time
}
