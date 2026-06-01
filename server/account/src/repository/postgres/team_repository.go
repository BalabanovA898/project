package postgres

import (
	"context"
	"database/sql"

	"github.com/BalabanovA898/project/account/src/domain"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team domain.Team) (string, error) {
	query := `
		INSERT INTO teams (name, owner_id, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`

	var id string
	if err := r.db.QueryRowContext(ctx, query, team.Name, team.OwnerID).Scan(&id); err != nil {
		return "", err
	}

	// Добавляем создателя как первого члена команды
	memberQuery := `
		INSERT INTO team_members (team_id, user_id, joined_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
	`
	if _, err := r.db.ExecContext(ctx, memberQuery, id, team.OwnerID); err != nil {
		return "", err
	}

	return id, nil
}

func (r *TeamRepository) GetByID(ctx context.Context, teamID string) (domain.Team, error) {
	query := `
		SELECT id, name, owner_id, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	var team domain.Team
	err := r.db.QueryRowContext(ctx, query, teamID).Scan(
		&team.ID, &team.Name, &team.OwnerID, &team.CreatedAt, &team.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Team{}, domain.ErrNotFound
		}
		return domain.Team{}, err
	}

	// Получаем членов команды
	members, err := r.getMembers(ctx, teamID)
	if err != nil {
		return domain.Team{}, err
	}
	team.Members = members

	return team, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, teamID, userID string) error {
	// Проверяем, что команда существует
	query := `SELECT id FROM teams WHERE id = $1`
	if err := r.db.QueryRowContext(ctx, query, teamID).Scan(&teamID); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return err
	}

	// Проверяем, что пользователь не уже в команде
	memberQuery := `SELECT user_id FROM team_members WHERE team_id = $1 AND user_id = $2`
	var existingUserID string
	err := r.db.QueryRowContext(ctx, memberQuery, teamID, userID).Scan(&existingUserID)
	if err == nil {
		return domain.ErrConflict // Уже в команде
	} else if err != sql.ErrNoRows {
		return err
	}

	// Добавляем члена
	insertQuery := `
		INSERT INTO team_members (team_id, user_id, joined_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
	`
	if _, err := r.db.ExecContext(ctx, insertQuery, teamID, userID); err != nil {
		return err
	}

	return nil
}

func (r *TeamRepository) RemoveMember(ctx context.Context, teamID, userID string) error {
	query := `
		DELETE FROM team_members
		WHERE team_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, teamID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TeamRepository) getMembers(ctx context.Context, teamID string) ([]string, error) {
	query := `
		SELECT user_id FROM team_members WHERE team_id = $1 ORDER BY joined_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}

	return members, rows.Err()
}
