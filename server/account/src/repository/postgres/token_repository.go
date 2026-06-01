package postgres

import (
	"context"
	"database/sql"

	"github.com/BalabanovA898/project/account/src/domain"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, token.Token, token.UserID, token.ExpiresAt)
	return err
}

func (r *RefreshTokenRepository) Get(ctx context.Context, token string) (domain.RefreshToken, error) {
	var result domain.RefreshToken
	query := `SELECT token, user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&result.Token, &result.UserID, &result.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.RefreshToken{}, domain.ErrNotFound
		}
		return domain.RefreshToken{}, err
	}

	return result, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
