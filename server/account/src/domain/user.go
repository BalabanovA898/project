package domain

import "time"

type Preferences struct {
	Priority        string `json:"priority,omitempty"`
	UnavailableDays string `json:"unavailableDays,omitempty"`
}

type User struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	Username     string       `json:"username"`
	PasswordHash string       `json:"-"`
	Disabled     bool         `json:"disabled"`
	Preferences  *Preferences `json:"preferences,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
