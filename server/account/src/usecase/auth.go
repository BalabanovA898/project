package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"time"
	"unicode"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/BalabanovA898/project/account/src/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	users      *postgres.UserRepository
	tokens     *postgres.RefreshTokenRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthUseCase(
	users *postgres.UserRepository,
	tokens *postgres.RefreshTokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		users:      users,
		tokens:     tokens,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, email, password, username string) (domain.User, domain.TokenResponse, error) {
	// Валидация email
	if err := validateEmail(email); err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	// Валидация пароля
	if err := validatePassword(password); err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	if _, err := u.users.GetByEmail(ctx, email); err == nil {
		return domain.User{}, domain.TokenResponse{}, domain.ErrConflict
	} else if err != domain.ErrNotFound {
		return domain.User{}, domain.TokenResponse{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	user := domain.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashed),
		Disabled:     false,
		Preferences: &domain.Preferences{
			Priority:        "none",
			UnavailableDays: "",
		},
	}

	id, err := u.users.Create(ctx, user)
	if err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	user.ID = id

	tokens, err := u.issueTokens(ctx, user)
	if err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	return user, tokens, nil
}

func (u *AuthUseCase) Login(ctx context.Context, email, password string) (domain.User, domain.TokenResponse, error) {
	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.User{}, domain.TokenResponse{}, domain.ErrInvalidCredential
		}
		return domain.User{}, domain.TokenResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, domain.TokenResponse{}, domain.ErrInvalidCredential
	}

	tokens, err := u.issueTokens(ctx, user)
	if err != nil {
		return domain.User{}, domain.TokenResponse{}, err
	}

	return user, tokens, nil
}

func (u *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (domain.TokenResponse, error) {
	stored, err := u.tokens.Get(ctx, refreshToken)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.TokenResponse{}, domain.ErrUnauthorized
		}
		return domain.TokenResponse{}, err
	}

	if stored.ExpiresAt.Before(time.Now()) {
		_ = u.tokens.Delete(ctx, refreshToken)
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}

	user, err := u.users.GetByID(ctx, stored.UserID)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	_ = u.tokens.Delete(ctx, refreshToken)

	tokens, err := u.issueTokens(ctx, user)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	return tokens, nil
}

func (u *AuthUseCase) issueTokens(ctx context.Context, user domain.User) (domain.TokenResponse, error) {
	access, err := u.buildAccessToken(user)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	refresh, err := u.buildRefreshToken(ctx, user.ID)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (u *AuthUseCase) buildAccessToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(u.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}

func (u *AuthUseCase) buildRefreshToken(ctx context.Context, userID string) (string, error) {
	payload := make([]byte, 32)
	if _, err := rand.Read(payload); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(payload)
	expiresAt := time.Now().Add(u.refreshTTL)

	if err := u.tokens.Create(ctx, domain.RefreshToken{Token: token, UserID: userID, ExpiresAt: expiresAt}); err != nil {
		return "", err
	}

	return token, nil
}

// validateEmail проверяет формат email адреса
func validateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return domain.NewValidationError("email", "invalid email format")
	}
	return nil
}

// validatePassword проверяет требования к паролю:
// - длина больше 8 символов
// - содержит большую букву
// - содержит маленькую букву
// - содержит специальный символ
func validatePassword(password string) error {
	if len(password) < 8 {
		return domain.NewValidationError("password", "password must be at least 8 characters long")
	}

	hasUpperCase := false
	hasLowerCase := false
	hasSpecialChar := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			hasSpecialChar = true
		}
	}

	if !hasUpperCase {
		return domain.NewValidationError("password", "password must contain at least one uppercase letter")
	}

	if !hasLowerCase {
		return domain.NewValidationError("password", "password must contain at least one lowercase letter")
	}

	if !hasSpecialChar {
		return domain.NewValidationError("password", "password must contain at least one special character")
	}

	return nil
}
