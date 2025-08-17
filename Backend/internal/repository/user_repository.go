package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rafael-bit/whatz/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, avatar, status, role, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Avatar, user.Status, user.Role, user.Tags, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar usuário: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, username, email, avatar, status, role, tags, created_at, updated_at
		FROM users WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Status, &user.Role, &user.Tags, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, avatar, status, role, tags, created_at, updated_at
		FROM users WHERE username = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Status, &user.Role, &user.Tags, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, username, email, avatar, status, role, tags, created_at, updated_at
		FROM users ORDER BY username
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Status, &user.Role, &user.Tags, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear usuário: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE users SET status = ?, updated_at = ? WHERE id = ?
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do usuário: %v", err)
	}

	return nil
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %v", err)
	}

	return nil
}

func (r *UserRepository) UpdateTags(id, tags string) error {
	query := `
		UPDATE users SET tags = ?, updated_at = ? WHERE id = ?
	`

	_, err := r.db.Exec(query, tags, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar tags do usuário: %v", err)
	}

	return nil
}

func (r *UserRepository) UpdateRole(id, role string) error {
	query := `
		UPDATE users SET role = ?, updated_at = ? WHERE id = ?
	`

	_, err := r.db.Exec(query, role, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar role do usuário: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByRole(role string) ([]*models.User, error) {
	query := `
		SELECT id, username, email, avatar, status, role, tags, created_at, updated_at
		FROM users WHERE role = ? ORDER BY username
	`

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários por role: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Status, &user.Role, &user.Tags, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear usuário: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}
