package repository

import (
	"database/sql"
	"fmt"

	"github.com/rafael-bit/whatz/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	query := `
		INSERT INTO messages (id, content, user_id, username, avatar, type, room_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, message.ID, message.Content, message.UserID, message.Username, message.Avatar, message.Type, message.RoomID, message.CreatedAt, message.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar mensagem: %v", err)
	}

	return nil
}

func (r *MessageRepository) GetByID(id string) (*models.Message, error) {
	query := `
		SELECT id, content, user_id, username, avatar, type, room_id, created_at, updated_at
		FROM messages WHERE id = ?
	`

	message := &models.Message{}
	err := r.db.QueryRow(query, id).Scan(
		&message.ID, &message.Content, &message.UserID, &message.Username, &message.Avatar, &message.Type, &message.RoomID, &message.CreatedAt, &message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar mensagem: %v", err)
	}

	return message, nil
}

func (r *MessageRepository) GetByRoom(roomID string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, content, user_id, username, avatar, type, room_id, created_at, updated_at
		FROM messages WHERE room_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens da sala: %v", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.Content, &message.UserID, &message.Username, &message.Avatar, &message.Type, &message.RoomID, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear mensagem: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) GetRecentMessages(roomID string, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, content, user_id, username, avatar, type, room_id, created_at, updated_at
		FROM messages WHERE room_id = ? ORDER BY created_at DESC LIMIT ?
	`

	rows, err := r.db.Query(query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens recentes: %v", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.Content, &message.UserID, &message.Username, &message.Avatar, &message.Type, &message.RoomID, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear mensagem: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) GetByUser(userID string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, content, user_id, username, avatar, type, room_id, created_at, updated_at
		FROM messages WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens do usuário: %v", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.Content, &message.UserID, &message.Username, &message.Avatar, &message.Type, &message.RoomID, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear mensagem: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar mensagem: %v", err)
	}

	return nil
}

func (r *MessageRepository) GetMessageCount(roomID string) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE room_id = ?`

	var count int
	err := r.db.QueryRow(query, roomID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar mensagens: %v", err)
	}

	return count, nil
}
