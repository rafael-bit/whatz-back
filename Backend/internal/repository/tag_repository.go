package repository

import (
	"database/sql"
	"fmt"

	"github.com/rafael-bit/whatz/internal/models"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	query := `
		INSERT INTO tags (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, tag.ID, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar tag: %v", err)
	}

	return nil
}

func (r *TagRepository) GetAll() ([]*models.Tag, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM tags ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tags: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear tag: %v", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) GetByName(name string) (*models.Tag, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM tags WHERE name = ?
	`

	tag := &models.Tag{}
	err := r.db.QueryRow(query, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar tag: %v", err)
	}

	return tag, nil
}

func (r *TagRepository) Delete(id string) error {
	query := `DELETE FROM tags WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar tag: %v", err)
	}

	return nil
}
