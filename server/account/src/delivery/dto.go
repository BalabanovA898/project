package delivery

import (
	"github.com/BalabanovA898/project/account/src/domain"
)

type UserPreviewDTO struct {
	ID          string                 `json:"id"`
	Username    string                 `json:"username"`
	Preferences *domain.Preferences    `json:"preferences,omitempty"`
}

func toUserPreviewDTO(user domain.User) UserPreviewDTO {
	return UserPreviewDTO{
		ID:          user.ID,
		Username:    user.Username,
		Preferences: user.Preferences,
	}
}
