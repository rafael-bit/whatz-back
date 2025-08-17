package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rafael-bit/whatz/internal/models"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *models.Room) error {
	query := `
		INSERT INTO rooms (id, name, description, type, access_tags, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, room.ID, room.Name, room.Description, room.Type, room.AccessTags, room.CreatedBy, room.CreatedAt, room.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar sala: %v", err)
	}

	return nil
}

func (r *RoomRepository) GetByID(id string) (*models.Room, error) {
	query := `
		SELECT id, name, description, type, access_tags, created_by, created_at, updated_at
		FROM rooms WHERE id = ?
	`

	room := &models.Room{}
	err := r.db.QueryRow(query, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.Type, &room.AccessTags, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar sala: %v", err)
	}

	return room, nil
}

func (r *RoomRepository) GetAll() ([]*models.Room, error) {
	query := `
		SELECT id, name, description, type, access_tags, created_by, created_at, updated_at
		FROM rooms ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salas: %v", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.Type, &room.AccessTags, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear sala: %v", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) GetPublicRooms() ([]*models.Room, error) {
	query := `
		SELECT id, name, description, type, access_tags, created_by, created_at, updated_at
		FROM rooms WHERE type = 'public' ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salas públicas: %v", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.Type, &room.AccessTags, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear sala: %v", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) GetByCreator(createdBy string) ([]*models.Room, error) {
	query := `
		SELECT id, name, description, type, access_tags, created_by, created_at, updated_at
		FROM rooms WHERE created_by = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, createdBy)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salas do criador: %v", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.Type, &room.AccessTags, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear sala: %v", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) Update(room *models.Room) error {
	query := `
		UPDATE rooms SET name = ?, description = ?, type = ?, access_tags = ?, updated_at = ? WHERE id = ?
	`

	_, err := r.db.Exec(query, room.Name, room.Description, room.Type, room.AccessTags, time.Now(), room.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar sala: %v", err)
	}

	return nil
}

func (r *RoomRepository) Delete(id string) error {
	query := `DELETE FROM rooms WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar sala: %v", err)
	}

	return nil
}

func (r *RoomRepository) GetRoomsByAccessTags(userTags []string) ([]*models.Room, error) {
	allRooms, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var accessibleRooms []*models.Room
	for _, room := range allRooms {
		if room.Type == "public" {
			accessibleRooms = append(accessibleRooms, room)
			continue
		}

		var roomTags []string
		if err := json.Unmarshal([]byte(room.AccessTags), &roomTags); err != nil {
			continue
		}

		if r.hasRequiredTags(roomTags, userTags) {
			accessibleRooms = append(accessibleRooms, room)
		}
	}

	return accessibleRooms, nil
}

func (r *RoomRepository) hasRequiredTags(roomTags, userTags []string) bool {
	if len(roomTags) == 0 {
		return true
	}

	for _, roomTag := range roomTags {
		for _, userTag := range userTags {
			if roomTag == userTag {
				return true
			}
		}
	}

	return false
}
