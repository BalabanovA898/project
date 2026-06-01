package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (string, error) {
	preferences, err := json.Marshal(user.Preferences)
	if err != nil {
		return "", err
	}

	var id string
	query := `INSERT INTO users (email, username, password_hash, disabled, preferences) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = r.db.QueryRowContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.Disabled, preferences).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return "", domain.ErrConflict
		}
		return "", err
	}

	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	var rawPreferences []byte
	query := `SELECT id, email, username, password_hash, disabled, preferences, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Disabled,
		&rawPreferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	if len(rawPreferences) > 0 {
		var prefs domain.Preferences
		if err := json.Unmarshal(rawPreferences, &prefs); err == nil {
			user.Preferences = &prefs
		}
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	var rawPreferences []byte
	query := `SELECT id, email, username, password_hash, disabled, preferences, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Disabled,
		&rawPreferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	if len(rawPreferences) > 0 {
		var prefs domain.Preferences
		if err := json.Unmarshal(rawPreferences, &prefs); err == nil {
			user.Preferences = &prefs
		}
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	preferences, err := json.Marshal(user.Preferences)
	if err != nil {
		return err
	}

	query := `UPDATE users SET username = $1, preferences = $2, updated_at = now() WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, user.Username, preferences, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
