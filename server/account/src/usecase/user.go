package usecase

import (
	"context"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/BalabanovA898/project/account/src/repository/postgres"
)

type UserUseCase struct {
	users *postgres.UserRepository
}

type UserPatch struct {
	Username    *string              `json:"username,omitempty"`
	Preferences *domain.Preferences `json:"preferences,omitempty"`
}

func NewUserUseCase(users *postgres.UserRepository) *UserUseCase {
	return &UserUseCase{users: users}
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return u.users.GetByID(ctx, id)
}

func (u *UserUseCase) UpdateUser(ctx context.Context, id string, patch UserPatch) (domain.User, error) {
	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	if patch.Username != nil {
		user.Username = *patch.Username
	}

	if patch.Preferences != nil {
		user.Preferences = patch.Preferences
	}

	if err := u.users.Update(ctx, user); err != nil {
		return domain.User{}, err
	}

	return u.users.GetByID(ctx, id)
}
