package repository

import (
	"database/sql"
	"time"

	"inventory-juanfe/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, is_active, last_login, created_at, updated_at
		FROM users
		WHERE email = $1`

	var u models.User
	err := r.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.IsActive,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, is_active, last_login, created_at, updated_at
		FROM users
		WHERE id = $1`

	var u models.User
	err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.IsActive,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) UpdateLastLogin(id string, t time.Time) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_login = $1 WHERE id = $2`,
		t, id,
	)
	return err
}
